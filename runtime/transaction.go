package runtime

import (
	"fmt"
	"github.com/fd/simplex/cas"
	"github.com/fd/simplex/cas/btree"
	"time"
)

type (
	Transaction struct {
		env     *Environment
		changes []*Change
		errors  []interface{}
		pool    *worker_pool_t

		// parent transaction
		Parent cas.Addr
		Tables *btree.Tree
	}

	ChangeKind uint

	Change struct {
		Kind  ChangeKind
		Table string
		Key   interface{}
		Elt   interface{}
	}
)

const (
	SET ChangeKind = iota
	UNSET
)

func (env *Environment) Transaction() *Transaction {
	txn := &Transaction{env: env}

	if tip_sha, ok := env.GetCurrentTransaction(); ok {
		var parent *Transaction

		ok := env.store.Get(tip_sha, &parent)
		if !ok {
			panic("corrupted data store")
		}

		// copy the *InternalTable structure
		txn.Tables = parent.Tables
		txn.Parent = tip_sha
	}

	if txn.Tables == nil {
		txn.Tables = &InternalTable{
			Name: "_tables",
		}
	}

	txn.Tables.txn = txn
	txn.Tables.setup()

	return txn
}

func (txn *Transaction) Set(table Table, key interface{}, val interface{}) {
	change := &Change{SET, table.TableId(), key, val}
	txn.changes = append(txn.changes, change)
}

func (txn *Transaction) Unset(table Table, key interface{}) {
	change := &Change{UNSET, table.TableId(), key, nil}
	txn.changes = append(txn.changes, change)
}

func (txn *Transaction) Commit() {
	var txn_sha storage.SHA
	{
		now := time.Now()
		defer func() {
			duration := time.Now().Sub(now)
			fmt.Printf("[sha: %s, duration: %s]\n", txn_sha, duration)
		}()
	}

	// wait for prev txn to resolve

	pool := &worker_pool_t{}
	txn.pool = pool
	events := pool.run()

	for _, t := range txn.env.terminals {
		pool.schedule(txn, t)
	}

	for event := range events {
		// handle events
		fmt.Printf("Ev (%T): %+v\n", event, event)
	}

	// commit the _tables table
	txn.Tables.Commit()
	txn_sha = txn.env.store.Set(&txn)
	txn.env.SetCurrentTransaction(txn_sha, txn.Parent)
}

func (txn *Transaction) GetTable(name string) *btree.Tree {
	_, elt_addr, err := txn.Tables.Get(cas.Collate(name))
	if cas.IsNotFound(err) {
		return btree.New(txn.env.Store)
	}
	if err != nil {
		panic("runtime: " + err.Error())
	}

	return txn.env.GetTable(elt_addr)
}

func (txn *Transaction) CommitTable(name string, tree *btree.Tree) (prev, curr cas.Addr) {
	var (
		key_coll      []byte
		key_addr      cas.Addr
		elt_addr      cas.Addr
		prev_elt_addr cas.Addr
		err           error
	)

	key_coll = cas.Collate(name)

	key_addr, err = cas.Encode(txn.env.Store, name)
	if err != nil {
		panic("runtime: " + err.Error())
	}

	elt_addr, err = tree.Commit()
	if err != nil {
		panic("runtime: " + err.Error())
	}

	_, prev_elt_addr, err = txn.Tables.Set(cas.Collate(name), elt_addr)
	if err != nil {
		panic("runtime: " + err.Error())
	}

	return prev_elt_addr, elt_addr
}

func (txn *Transaction) Resolve(def Deferred) <-chan Event {
	if txn.pool == nil {
		panic("transcation has no running worker pool")
	}

	worker := txn.pool.schedule(txn, def)
	return worker.subscribe()
}

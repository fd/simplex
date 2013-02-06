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
		tables  *btree.Tree
		errors  []interface{}
		pool    *worker_pool_t

		// parent transaction
		Parent cas.Addr
		Tables cas.Addr
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

	txn_addr, err := env.GetCurrentTransaction()
	if err != nil {
		panic("runtime: " + err.Error())
	}
	if txn_addr != nil {
		var parent *Transaction

		err := cas.Decode(env.Store, txn_addr, &parent)
		if err != nil {
			panic("runtime: " + err.Error())
		}

		// copy the *InternalTable structure
		txn.Tables = parent.Tables
		txn.Parent = txn_addr
	}

	if txn.Tables == nil {
		txn.tables = btree.New(txn.env.Store)
	} else {
		tables, err := btree.Open(env.Store, txn.Tables)
		if err != nil {
			panic("runtime: " + err.Error())
		}
		txn.tables = tables
	}

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
	var txn_addr cas.Addr
	{
		now := time.Now()
		defer func() {
			duration := time.Now().Sub(now)
			fmt.Printf("[sha: %s, duration: %s]\n", txn_addr, duration)
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
	tables_addr, err := txn.tables.Commit()
	if err != nil {
		panic("runtime: " + err.Error())
	}
	txn.Tables = tables_addr

	// overflow trigger is 0; we always write a transaction
	txn_addr, err = cas.Encode(txn.env.Store, &txn, 0)
	if err != nil {
		panic("runtime: " + err.Error())
	}

	err = txn.env.SetCurrentTransaction(txn_addr, txn.Parent)
	if err != nil {
		panic("runtime: " + err.Error())
	}
}

func (txn *Transaction) GetTable(name string) *btree.Tree {
	_, elt_addr, err := txn.tables.Get(cas.Collate(name))
	if err != nil {
		panic("runtime: " + err.Error())
	}

	if elt_addr == nil {
		return btree.New(txn.env.Store)
	}

	tree, err := btree.Open(txn.env.Store, elt_addr)
	if err != nil {
		panic("runtime: " + err.Error())
	}

	return tree
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

	key_addr, err = cas.Encode(txn.env.Store, name, -1)
	if err != nil {
		panic("runtime: " + err.Error())
	}

	elt_addr, err = tree.Commit()
	if err != nil {
		panic("runtime: " + err.Error())
	}

	prev_elt_addr, err = txn.tables.Set(key_coll, key_addr, elt_addr)
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

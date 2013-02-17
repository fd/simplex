package runtime

import (
	"fmt"
	"runtime/debug"
	"simplex.sh/cas"
	"simplex.sh/cas/btree"
	"sync"
	"time"
)

type (
	Transaction struct {
		env          *Environment
		changes      []*Change
		tables       *btree.Tree
		errors       []interface{}
		broadcasters map[string]broadcaster
		mutex        sync.Mutex

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

func NewTransaction(env *Environment) *Transaction {
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

	promises := make([]broadcaster, len(txn.env.terminals))

	for i, t := range txn.env.terminals {
		promises[i] = txn.GoResolve(t)
	}

	for _, b := range promises {
		b <- <-b
	}

	for _, err := range txn.errors {
		fmt.Printf("[E] %s\n", err)
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

func (txn *Transaction) Resolve(r Resolver) IChange {
	b := txn.GoResolve(r)
	c := <-b
	b <- c
	return c
}

func (txn *Transaction) GoResolve(r Resolver) broadcaster {
	txn.mutex.Lock()
	defer txn.mutex.Unlock()

	if txn.broadcasters == nil {
		txn.broadcasters = make(map[string]broadcaster, 50)
	}

	b := txn.broadcasters[r.DeferredId()]
	if b != nil {
		return b
	}

	b = make(chan IChange, 1)
	txn.broadcasters[r.DeferredId()] = b
	go func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				txn.errors = append(txn.errors, err)
				b <- IChange{Err: err, Stack: debug.Stack()}
			} else {
				txn.errors = append(txn.errors, fmt.Errorf("panic: %+v", e))
				b <- IChange{Err: fmt.Errorf("panic: %+v", e), Stack: debug.Stack()}
			}
		}

		b <- r.Resolve(txn)
	}()

	return b
}

func (txn *Transaction) Store() cas.Store { return txn.env.Store }

type broadcaster chan IChange

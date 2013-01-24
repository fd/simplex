package runtime

import (
	"fmt"
	"github.com/fd/simplex/data/storage"
)

type (
	Transaction struct {
		env     *Environment
		changes []*Change
		errors  []interface{}

		// parent transaction
		Parent storage.SHA
		Tables *InternalTable
	}

	ChangeKind uint

	Change struct {
		Kind  ChangeKind
		Table string
		Key   interface{}
		Value interface{}
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
	// wait for prev txn to resolve

	pool := &worker_pool_t{}
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
	txn_sha := txn.env.store.Set(&txn)
	txn.env.SetCurrentTransaction(txn_sha, txn.Parent)
}

func (txn *Transaction) GetTable(name string) *InternalTable {
	var table *InternalTable

	ok := txn.Tables.Get(name, &table)
	if !ok {
		table = &InternalTable{
			txn:  txn,
			Name: name,
		}
		table.setup()
		return table
	}

	table.txn = txn
	table.setup()
	return table
}

func (txn *Transaction) Resolve(def ...Deferred) <-chan Event {
	panic("not implemented")
}

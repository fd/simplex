package runtime

import (
	"fmt"
)

type (
	Transaction struct {
		env     *Environment
		changes []*Change

		errors []interface{}
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
	return &Transaction{
		env: env,
	}
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
}

func (txn *Transaction) GetTable(name string) *InternalTable {
	var table *InternalTable

	table_ref, ok := txn.tables.Get(name)
	if !ok {
		table = &InternalTable{txn: txn}
		return table
	}

	ok = txn.storage.Get(table_ref, &table)
	if !ok {
		table = &InternalTable{txn: txn}
		return table
	}

	return table
}

func (txn *Transaction) Resolve(def ...Deferred) <-chan Event {
	panic("not implemented")
}

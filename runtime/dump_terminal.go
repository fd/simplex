package runtime

import (
	"fmt"
	"reflect"
)

func Dump(view IndexedView) {
	Env.RegisterTerminal(&dump_terminal{view})
}

type dump_terminal struct {
	view IndexedView
}

func (t *dump_terminal) DeferredId() string {
	return "dump(" + t.view.DeferredId() + ")"
}

func (t *dump_terminal) Resolve(txn *Transaction, events chan<- Event) {
	i_events := txn.Resolve(t.view)

	for e := range i_events {
		event, ok := e.(*EvConsistent)
		if !ok {
			continue
		}

		table := event.GetTableB(txn)
		iter := table.Iter()

		for {
			sha, done := iter.Next()
			if done {
				break
			}

			var (
				kv    KeyValue
				value reflect.Value
			)
			if !txn.env.store.Get(sha, &kv) {
				panic("corrupt")
			}

			value = reflect.New(t.view.EltType())
			if !txn.env.store.GetValue(kv.ValueSha, value) {
				panic("corrupt")
			}

			fmt.Printf("V: `%s` %+v\n", kv.KeyCompare, value.Interface())
		}
	}
}

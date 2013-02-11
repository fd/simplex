package runtime

import (
	"fmt"
	"github.com/fd/simplex/cas"
	"github.com/fd/simplex/cas/btree"
	"github.com/fd/simplex/runtime/event"
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

func (t *dump_terminal) Resolve(txn *Transaction, events chan<- event.Event) {
	src_events := txn.Resolve(t.view)

	for e := range src_events.C {
		// propagate error events
		if err, ok := e.(event.Error); ok {
			events <- err
			continue
		}

		event, ok := e.(*ConsistentTable)
		if !ok {
			continue
		}

		var (
			table   = event.GetTableB(txn)
			iter    = table.Iter()
			keyed   bool
			key_typ reflect.Type
		)

		if kv, ok := t.view.(KeyedView); ok {
			keyed = true
			key_typ = kv.KeyType()
		}

		for {
			key_addr, elt_addr, err := iter.Next()
			if err == btree.EOI {
				err = nil
				break
			}
			if err != nil {
				panic("runtime: " + err.Error())
			}

			var (
				key reflect.Value
				elt reflect.Value
			)

			if keyed {
				key = reflect.New(key_typ)
				err = cas.DecodeValue(txn.env.Store, key_addr, key)
				if err != nil {
					panic("runtime: " + err.Error())
				}
			}

			elt = reflect.New(t.view.EltType())
			err = cas.DecodeValue(txn.env.Store, elt_addr, elt)
			if err != nil {
				panic("runtime: " + err.Error())
			}

			if keyed {
				fmt.Printf("V: %+v %+v\n", key.Interface(), elt.Interface())
			} else {
				fmt.Printf("V: %+v\n", elt.Interface())
			}
		}
	}
}

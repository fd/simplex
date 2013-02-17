package runtime

import (
	"fmt"
	"reflect"
	"simplex.sh/cas"
	"simplex.sh/cas/btree"
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

func (t *dump_terminal) Resolve(state *Transaction) IChange {
	var (
		i_change = state.Resolve(t.view)
	)

	if i_change.Type() == ChangeRemove {
		return IChange{}
	}

	var (
		table   = GetTable(state.Store(), i_change.B)
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
			err = cas.DecodeValue(state.Store(), key_addr, key)
			if err != nil {
				panic("runtime: " + err.Error())
			}
		}

		elt = reflect.New(t.view.EltType())
		err = cas.DecodeValue(state.Store(), elt_addr, elt)
		if err != nil {
			panic("runtime: " + err.Error())
		}

		if keyed {
			fmt.Printf("V: %+v %+v\n", key.Interface(), elt.Interface())
		} else {
			fmt.Printf("V: %+v\n", elt.Interface())
		}
	}

	return IChange{}
}

package runtime

import (
	"bytes"
	"github.com/fd/simplex/cas"
	"github.com/fd/simplex/runtime/event"
)

func (op *group_op) Resolve(txn *Transaction, events chan<- event.Event) {
	var (
		src_events = txn.Resolve(op.src)
		table      = txn.GetTable(op.name)
	)

	for e := range src_events.C {
		// propagate error events
		if err, ok := e.(event.Error); ok {
			events <- err
			continue
		}

		i_change, ok := e.(*ChangedMember)
		if !ok {
			continue
		}

		var (
			coll_key_a []byte
			coll_key_b []byte
			key_b      interface{}
		)

		// calculate collated group key for a and b
		if i_change.a != nil {
			group_key := op.fun(&Context{txn}, i_change.a)
			coll_key_a = cas.Collate(group_key)
		}
		if i_change.b != nil {
			key_b = op.fun(&Context{txn}, i_change.b)
			coll_key_b = cas.Collate(key_b)
		}

		// skip when they are equal
		if bytes.Compare(coll_key_a, coll_key_b) == 0 {
			continue
		}

		// remove old entry
		if i_change.a != nil {
			key_addr, elt_addr, err := table.Del(coll_key_a)
			if err != nil {
				panic("runtime: " + err.Error())
			}
			events <- &ChangedMember{op.name, coll_key_a, key_addr, elt_addr, nil}
		}

		// add new entry
		if i_change.b != nil {
			key_addr, err := cas.Encode(txn.env.Store, key_b, -1)
			if err != nil {
				panic("runtime: " + err.Error())
			}

			prev_elt_addr, err := table.Set(coll_key_a, key_addr, i_change.b)
			if err != nil {
				panic("runtime: " + err.Error())
			}

			events <- &ChangedMember{op.name, coll_key_b, key_addr, prev_elt_addr, i_change.b}
		}
	}

	tab_addr_a, tab_addr_b := txn.CommitTable(op.name, table)
	events <- &ConsistentTable{op.name, tab_addr_a, tab_addr_b}
}

package runtime

import (
	"bytes"
	"github.com/fd/simplex/cas"
)

func (op *sort_op) Resolve(txn *Transaction, events chan<- Event) {
	var (
		src_events = txn.Resolve(op.src)
		table      = txn.GetTable(op.name)
	)

	for event := range src_events {
		i_change, ok := event.(*ev_CHANGE)
		if !ok {
			continue
		}

		var (
			coll_key_a []byte
			coll_key_b []byte
			key_b      []interface{}
		)

		// calculate collated sort key for a and b
		if i_change.a != nil {
			sort_key := op.fun(&Context{txn}, i_change.a)
			coll_key_a = cas.Collate([]interface{}{sort_key, i_change.collated_key})
		}
		if i_change.b != nil {
			sort_key := op.fun(&Context{txn}, i_change.b)
			key_b = []interface{}{sort_key, i_change.collated_key}
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
			events <- &ev_CHANGE{op.name, coll_key_a, key_addr, elt_addr, nil}
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

			events <- &ev_CHANGE{op.name, coll_key_b, key_addr, prev_elt_addr, i_change.b}
		}
	}

	tab_addr_a, tab_addr_b := txn.CommitTable(op.name, table)
	events <- &EvConsistent{op.name, tab_addr_a, tab_addr_b}
}

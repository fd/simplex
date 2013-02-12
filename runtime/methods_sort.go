package runtime

import (
	"bytes"
	"simplex.sh/cas"
	"simplex.sh/runtime/event"
	"simplex.sh/runtime/promise"
)

func (op *sort_op) Resolve(state promise.State, events chan<- event.Event) {
	var (
		src_events = state.Resolve(op.src)
		table      = state.GetTable(op.name)
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
			key_b      []interface{}
		)

		// calculate collated sort key for a and b
		if i_change.a != nil {
			sort_key := op.fun(&Context{state.Store()}, i_change.a)
			coll_key_a = cas.Collate([]interface{}{sort_key, i_change.collated_key})
		}
		if i_change.b != nil {
			sort_key := op.fun(&Context{state.Store()}, i_change.b)
			key_b = []interface{}{sort_key, i_change.collated_key}
			coll_key_b = cas.Collate(key_b)
		}

		// forward when the keys are equal
		if bytes.Compare(coll_key_a, coll_key_b) == 0 {
			key_addr, err := cas.Encode(state.Store(), key_b, -1)
			if err != nil {
				panic("runtime: " + err.Error())
			}

			events <- &ChangedMember{op.name, coll_key_b, key_addr, i_change.a, i_change.b}
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
			key_addr, err := cas.Encode(state.Store(), key_b, -1)
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

	tab_addr_a, tab_addr_b := state.CommitTable(op.name, table)
	events <- &ConsistentTable{op.name, tab_addr_a, tab_addr_b}
}

package runtime

import (
	"bytes"
	"simplex.sh/cas"
	"simplex.sh/runtime/event"
	"simplex.sh/runtime/promise"
)

func (op *group_op) Resolve(state promise.State, events chan<- event.Event) {
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
			key_b      interface{}
		)

		// calculate collated group key for a and b
		if i_change.a != nil {
			group_key := op.fun(&Context{state.Store()}, i_change.a)
			coll_key_a = cas.Collate(group_key)
		}
		if i_change.b != nil {
			key_b = op.fun(&Context{state.Store()}, i_change.b)
			coll_key_b = cas.Collate(key_b)
		}

		// propagate event
		// - to sub table at coll_key_b
		// - to groups table
		if bytes.Compare(coll_key_a, coll_key_b) == 0 {
			continue
		}

		// remove old entry from sub table
		// add new entry to sub table (while potentially adding new subtables)
	}

	// remove empty sub tables

	tab_addr_a, tab_addr_b := state.CommitTable(op.name, table)
	events <- &ConsistentTable{op.name, tab_addr_a, tab_addr_b}
}

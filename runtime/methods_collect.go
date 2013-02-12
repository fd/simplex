package runtime

import (
	"simplex.sh/cas"
	"simplex.sh/runtime/event"
	"simplex.sh/runtime/promise"
)

func (op *collect_op) Resolve(state promise.State, events chan<- event.Event) {
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

		// removed
		if i_change.b == nil {
			prev_key_addr, prev_elt_addr, err := table.Del(i_change.collated_key)
			if err != nil {
				panic("runtime: " + err.Error())
			}

			if prev_key_addr != nil && prev_elt_addr != nil {
				events <- &ChangedMember{op.name, i_change.collated_key, prev_key_addr, prev_elt_addr, nil}
			}

			continue
		}

		{ // added or updated
			curr_elt_addr := op.fun(&Context{state.Store()}, i_change.b)

			prev_elt_addr, err := table.Set(i_change.collated_key, i_change.key, curr_elt_addr)
			if err != nil {
				panic("runtime: " + err.Error())
			}
			if cas.CompareAddr(prev_elt_addr, curr_elt_addr) != 0 {
				events <- &ChangedMember{op.name, i_change.collated_key, i_change.key, prev_elt_addr, curr_elt_addr}
			}
		}
	}

	tab_addr_a, tab_addr_b := state.CommitTable(op.name, table)
	events <- &ConsistentTable{op.name, tab_addr_a, tab_addr_b}
}

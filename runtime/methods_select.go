package runtime

import (
	"simplex.sh/cas"
	"simplex.sh/runtime/event"
	"simplex.sh/runtime/promise"
)

func (op *select_op) Resolve(state promise.State, events chan<- event.Event) {
	var (
		src_events = state.Resolve(op.src)
		fun        = op.fun
	)

	apply_select_reject_filter(op.name, fun, true, src_events, events, state)
}

func (op *reject_op) Resolve(state promise.State, events chan<- event.Event) {
	var (
		src_events = state.Resolve(op.src)
		fun        = select_func(op.fun)
	)

	apply_select_reject_filter(op.name, fun, false, src_events, events, state)
}

func apply_select_reject_filter(op_name string, op_fun select_func,
	expected bool, src_events *event.Subscription, dst_events chan<- event.Event,
	state promise.State) {

	var (
		table = state.GetTable(op_name)
	)

	for e := range src_events.C {
		// propagate error events
		if err, ok := e.(event.Error); ok {
			dst_events <- err
			continue
		}

		i_change, ok := e.(*ChangedMember)
		if !ok {
			continue
		}

		var (
			o_change = &ChangedMember{op_name, i_change.collated_key, i_change.key, i_change.a, i_change.b}
		)

		if o_change.a != nil {
			if op_fun(&Context{state.Store()}, o_change.a) != expected {
				o_change.a = nil
			}
		}

		if o_change.b != nil {
			if op_fun(&Context{state.Store()}, o_change.b) != expected {
				o_change.b = nil
			}
		}

		// ignore unchanged data
		if o_change.a == nil && o_change.b == nil {
			continue
		}

		if o_change.a != nil {
			// remove kv from table
			_, prev, err := table.Del(o_change.collated_key)
			if err != nil {
				panic("runtime: " + err.Error())
			}
			if prev != nil {
				o_change.a = nil
			}
		}

		if o_change.b != nil {
			// insert kv into table
			prev, err := table.Set(o_change.collated_key, o_change.key, o_change.b)
			if err != nil {
				panic("runtime: " + err.Error())
			}
			if cas.CompareAddr(prev, o_change.b) == 0 {
				o_change.b = nil
			}
		}

		// ignore unchanged data
		if o_change.a == nil && o_change.b == nil {
			continue
		}

		dst_events <- o_change
	}

	tab_addr_a, tab_addr_b := state.CommitTable(op_name, table)
	dst_events <- &ConsistentTable{op_name, tab_addr_a, tab_addr_b}
}

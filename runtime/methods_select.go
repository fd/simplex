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

		i_change, ok := e.(*Changed)
		if !ok {
			continue
		}

		if i_change.Depth() != 1 {
			continue
		}

		var (
			o_change = &Changed{
				[][]byte{
					[]byte(op_name),
					i_change.IdAt(-1),
				},
				i_change.Key,
				i_change.A,
				i_change.B,
			}
		)

		if o_change.A != nil {
			if op_fun(&Context{state.Store()}, o_change.A) != expected {
				o_change.A = nil
			}
		}

		if o_change.B != nil {
			if op_fun(&Context{state.Store()}, o_change.B) != expected {
				o_change.B = nil
			}
		}

		// forward changes
		if o_change.A == nil && o_change.B == nil {
			// TODO update table
			dst_events <- o_change
			continue
		}

		if o_change.A != nil {
			// remove kv from table
			_, prev, err := table.Del(o_change.IdAt(-1))
			if err != nil {
				panic("runtime: " + err.Error())
			}
			o_change.A = prev
		}

		if o_change.B != nil {
			// insert kv into table
			prev, err := table.Set(o_change.IdAt(-1), o_change.Key, o_change.B)
			if err != nil {
				panic("runtime: " + err.Error())
			}
			if cas.CompareAddr(prev, o_change.B) == 0 {
				o_change.B = nil
			}
		}

		// ignore unchanged data
		if o_change.A == nil && o_change.B == nil {
			continue
		}

		dst_events <- o_change
	}

	tab_addr_a, tab_addr_b := state.CommitTable(op_name, table)
	dst_events <- &Changed{
		[][]byte{
			[]byte(op_name),
		},
		nil,
		tab_addr_a,
		tab_addr_b,
	}
}

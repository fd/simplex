package runtime

import (
	"github.com/fd/simplex/cas"
)

func (op *select_op) Resolve(txn *Transaction, events chan<- Event) {
	var (
		src_events = txn.Resolve(op.src)
	)

	apply_select_reject_filter(op.name, op.fun, true, src_events, events, txn)
}

func (op *reject_op) Resolve(txn *Transaction, events chan<- Event) {
	var (
		src_events = txn.Resolve(op.src)
	)

	apply_select_reject_filter(op.name, select_func(op.fun), false, src_events, events, txn)
}

func apply_select_reject_filter(op_name string, op_fun select_func,
	expexted bool, src_events <-chan Event, dst_events chan<- Event,
	txn *Transaction) {

	var (
		table = txn.GetTable(op_name)
	)

	for event := range src_events {
		i_change, ok := event.(*ev_CHANGE)
		if !ok {
			continue
		}

		var (
			o_change = &ev_CHANGE{op_name, i_change.collated_key, i_change.key, i_change.a, i_change.b}
		)

		if o_change.a != nil {
			if op_fun(&Context{txn}, o_change.a) != expexted {
				o_change.a = nil
			}
		}

		if o_change.b != nil {
			if op_fun(&Context{txn}, o_change.b) != expexted {
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

	tab_addr_a, tab_addr_b := txn.CommitTable(op_name, table)
	dst_events <- &EvConsistent{op_name, tab_addr_a, tab_addr_b}
}

package runtime

import (
	"github.com/fd/simplex/cas"
)

func (op *collect_op) Resolve(txn *Transaction, events chan<- Event) {
	var (
		src_event = txn.Resolve(op.src)
		table     = txn.GetTable(op.name)
	)

	for event := range src_event {
		i_change, ok := event.(*ev_CHANGE)
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
				events <- &ev_CHANGE{op.name, i_change.collated_key, prev_key_addr, prev_elt_addr, nil}
			}

			continue
		}

		{ // added or updated
			curr_elt_addr := op.fun(&Context{txn}, i_change.b)

			prev_elt_addr, err := table.Set(i_change.collated_key, i_change.key, curr_elt_addr)
			if err != nil {
				panic("runtime: " + err.Error())
			}
			if cas.CompareAddr(prev_elt_addr, curr_elt_addr) != 0 {
				events <- &ev_CHANGE{op.name, i_change.collated_key, i_change.key, prev_elt_addr, curr_elt_addr}
			}
		}
	}

	tab_addr_a, tab_addr_b := txn.CommitTable(op.name, table)
	events <- &EvConsistent{op.name, tab_addr_a, tab_addr_b}
}

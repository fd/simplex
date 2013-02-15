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

		i_change, ok := e.(*Changed)
		if !ok {
			continue
		}

		if i_change.Depth() != 1 {
			continue
		}

		var (
			coll_key_a []byte
			coll_key_b []byte
			key_b      []interface{}
		)

		// calculate collated sort key for a and b
		if i_change.A != nil {
			sort_key := op.fun(&Context{state.Store()}, i_change.A)
			coll_key_a = cas.Collate([]interface{}{sort_key, i_change.IdAt(-1)})
		}
		if i_change.B != nil {
			sort_key := op.fun(&Context{state.Store()}, i_change.B)
			key_b = []interface{}{sort_key, i_change.IdAt(-1)}
			coll_key_b = cas.Collate(key_b)
		}

		// forward when the keys are equal
		if bytes.Compare(coll_key_a, coll_key_b) == 0 {
			key_addr, err := cas.Encode(state.Store(), key_b, -1)
			if err != nil {
				panic("runtime: " + err.Error())
			}

			events <- &Changed{
				[][]byte{
					[]byte(op.name),
					coll_key_b,
				},
				key_addr,
				i_change.A,
				i_change.B,
			}
			continue
		}

		// remove old entry
		if i_change.A != nil {
			key_addr, elt_addr, err := table.Del(coll_key_a)
			if err != nil {
				panic("runtime: " + err.Error())
			}

			events <- &Changed{
				[][]byte{
					[]byte(op.name),
					coll_key_a,
				},
				key_addr,
				elt_addr,
				nil,
			}
		}

		// add new entry
		if i_change.B != nil {
			key_addr, err := cas.Encode(state.Store(), key_b, -1)
			if err != nil {
				panic("runtime: " + err.Error())
			}

			prev_elt_addr, err := table.Set(coll_key_a, key_addr, i_change.B)
			if err != nil {
				panic("runtime: " + err.Error())
			}

			events <- &Changed{
				[][]byte{
					[]byte(op.name),
					coll_key_b,
				},
				key_addr,
				prev_elt_addr,
				i_change.B,
			}
		}
	}

	tab_addr_a, tab_addr_b := state.CommitTable(op.name, table)
	events <- &Changed{
		[][]byte{
			[]byte(op.name),
		},
		nil,
		tab_addr_a,
		tab_addr_b,
	}
}

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

		i_change, ok := e.(*Changed)
		if !ok {
			continue
		}

		if i_change.Depth() != 1 {
			continue
		}

		// removed
		if i_change.B == nil {
			prev_key_addr, prev_elt_addr, err := table.Del(i_change.IdAt(-1))
			if err != nil {
				panic("runtime: " + err.Error())
			}

			if prev_key_addr != nil && prev_elt_addr != nil {
				events <- &Changed{
					[][]byte{
						[]byte(op.name),
						i_change.IdAt(-1),
					},
					prev_key_addr,
					prev_elt_addr,
					nil,
				}
			}

			continue
		}

		{ // added or updated
			curr_elt_addr := op.fun(&Context{state.Store()}, i_change.B)

			prev_elt_addr, err := table.Set(i_change.IdAt(-1), i_change.Key, curr_elt_addr)
			if err != nil {
				panic("runtime: " + err.Error())
			}
			if cas.CompareAddr(prev_elt_addr, curr_elt_addr) != 0 {
				events <- &Changed{
					[][]byte{
						[]byte(op.name),
						i_change.IdAt(-1),
					},
					i_change.Key,
					prev_elt_addr,
					curr_elt_addr,
				}
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

package runtime

import (
	"bytes"
	"simplex.sh/cas"
)

func (op *sort_op) Resolve(state *Transaction) IChange {
	var (
		i_change = state.Resolve(op.src)
		o_change IChange
	)

	if i_change.Type() == ChangeNone {
		return o_change
	}

	var (
		table = state.GetTable(op.name)
	)

	for _, m := range i_change.MemberChanges {
		var (
			coll_key_a []byte
			coll_key_b []byte
			key_b      []interface{}
		)

		// calculate collated sort key for a and b
		if m.A != nil {
			sort_key := op.fun(&Context{state.Store()}, m.A)
			coll_key_a = cas.Collate([]interface{}{sort_key, m.CollatedKey})
		}
		if m.B != nil {
			sort_key := op.fun(&Context{state.Store()}, m.B)
			key_b = []interface{}{sort_key, m.CollatedKey}
			coll_key_b = cas.Collate(key_b)
		}

		// forward when the keys are equal
		if bytes.Compare(coll_key_a, coll_key_b) == 0 {
			key_addr, err := cas.Encode(state.Store(), key_b, -1)
			if err != nil {
				panic("runtime: " + err.Error())
			}

			prev_elt_addr, err := table.Set(coll_key_a, key_addr, m.B)
			if err != nil {
				panic("runtime: " + err.Error())
			}

			o_change.MemberChanged(coll_key_a, key_addr, IChange{A: prev_elt_addr, B: m.B})
			continue
		}

		// remove old entry
		if m.A != nil {
			_, prev_elt_addr, err := table.Del(coll_key_a)
			if err != nil {
				panic("runtime: " + err.Error())
			}

			o_change.MemberChanged(m.CollatedKey, m.Key, IChange{A: prev_elt_addr, B: nil})
		}

		// add new entry
		if m.B != nil {
			key_addr, err := cas.Encode(state.Store(), key_b, -1)
			if err != nil {
				panic("runtime: " + err.Error())
			}

			prev_elt_addr, err := table.Set(coll_key_b, key_addr, m.B)
			if err != nil {
				panic("runtime: " + err.Error())
			}

			o_change.MemberChanged(coll_key_b, key_addr, IChange{A: prev_elt_addr, B: m.B})
		}
	}

	o_change.A, o_change.B = state.CommitTable(op.name, table)
	return o_change
}

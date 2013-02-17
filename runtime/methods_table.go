package runtime

import (
	"simplex.sh/cas"
)

func (op *table_op) Resolve(state *Transaction) IChange {
	var (
		table    = state.GetTable(op.name)
		o_change IChange
	)

	for _, change := range state.changes {
		if change.Table != op.name {
			continue
		}

		switch change.Kind {
		case SET:
			var (
				key_coll      []byte
				key_addr      cas.Addr
				elt_addr      cas.Addr
				prev_elt_addr cas.Addr
				err           error
			)

			key_coll = cas.Collate(change.Key)

			key_addr, err = cas.Encode(state.Store(), change.Key, -1)
			if err != nil {
				panic("runtime: " + err.Error())
			}

			elt_addr, err = cas.Encode(state.Store(), change.Elt, -1)
			if err != nil {
				panic("runtime: " + err.Error())
			}

			prev_elt_addr, err = table.Set(key_coll, key_addr, elt_addr)
			if err != nil {
				panic("runtime: " + err.Error())
			}

			o_change.MemberChanged(key_coll, key_addr, IChange{A: prev_elt_addr, B: elt_addr})

		case UNSET:
			var (
				key_coll      []byte
				key_addr      cas.Addr
				prev_elt_addr cas.Addr
				err           error
			)

			key_coll = cas.Collate(change.Key)

			key_addr, prev_elt_addr, err = table.Del(key_coll)
			if err != nil {
				panic("runtime: " + err.Error())
			}

			o_change.MemberChanged(key_coll, key_addr, IChange{A: prev_elt_addr, B: nil})

		}
	}

	o_change.A, o_change.B = state.CommitTable(op.name, table)
	return o_change
}

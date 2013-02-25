package runtime

func (op *collect_op) Resolve(state *Transaction) IChange {
	var (
		i_change = state.Resolve(op.src)
		o_change = IChange{}
	)

	if i_change.Type() == ChangeNone {
		return o_change
	}

	var (
		table = state.GetTable(op.name)
	)

	for _, m := range i_change.MemberChanges {
		switch m.Type() {

		case ChangeRemove:
			_, prev_elt_addr, err := table.Del(m.CollatedKey)
			if err != nil {
				panic("runtime: " + err.Error())
			}

			o_change.MemberChanged(m.CollatedKey, m.Key, IChange{A: prev_elt_addr, B: nil})

		case ChangeUpdate, ChangeInsert:
			curr_elt_addr := op.fun(&Context{state.Store()}, m.B)

			prev_elt_addr, err := table.Set(m.CollatedKey, m.Key, curr_elt_addr)
			if err != nil {
				panic("runtime: " + err.Error())
			}

			o_change.MemberChanged(m.CollatedKey, m.Key, IChange{A: prev_elt_addr, B: curr_elt_addr})

		}
	}

	o_change.A, o_change.B = state.CommitTable(op.name, table)
	return o_change
}

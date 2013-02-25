package runtime

func (op *select_op) Resolve(state *Transaction) IChange {
	return apply_select_reject_filter(
		op.src,
		op.name,
		op.fun,
		true,
		state,
	)
}

func (op *reject_op) Resolve(state *Transaction) IChange {
	return apply_select_reject_filter(
		op.src,
		op.name,
		select_func(op.fun),
		false,
		state,
	)
}

func apply_select_reject_filter(r Resolver, op_name string, op_fun select_func,
	expected bool, state *Transaction) IChange {

	var (
		i_change = state.Resolve(r)
		o_change IChange
	)

	if i_change.Type() == ChangeNone {
		return o_change
	}

	var (
		table = state.GetTable(op_name)
	)

	for _, i_m := range i_change.MemberChanges {

		var (
			o_m = i_m.IChange
		)

		// was part of selection
		if i_m.A != nil {
			if op_fun(&Context{state.Store()}, i_m.A) != expected {
				o_m.A = nil
			}
		}

		// will be part of selection
		if i_m.B != nil {
			if op_fun(&Context{state.Store()}, i_m.B) != expected {
				o_m.B = nil
			}
		}

		if o_m.Type() == ChangeNone {
			continue
		}

		switch o_m.Type() {
		case ChangeRemove:
			_, prev_elt_addr, err := table.Del(i_m.CollatedKey)
			if err != nil {
				panic("runtime: " + err.Error())
			}

			o_m.A = prev_elt_addr
			o_change.MemberChanged(i_m.CollatedKey, i_m.Key, o_m)

		case ChangeUpdate, ChangeInsert:
			// insert kv into table
			prev_elt_addr, err := table.Set(i_m.CollatedKey, i_m.Key, o_m.B)
			if err != nil {
				panic("runtime: " + err.Error())
			}

			o_m.A = prev_elt_addr
			o_change.MemberChanged(i_m.CollatedKey, i_m.Key, o_m)

		}
	}

	o_change.A, o_change.B = state.CommitTable(op_name, table)
	return o_change
}

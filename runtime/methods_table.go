package runtime

import (
	"simplex.sh/cas"
	"simplex.sh/runtime/event"
	"simplex.sh/runtime/promise"
)

func (op *table_op) Resolve(state promise.State, events chan<- event.Event) {
	table := state.GetTable(op.name)

	transaction, ok := state.(*Transaction)
	if !ok {
		panic("Expected state to be a transaction.")
	}

	for _, change := range transaction.changes {
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

			if cas.CompareAddr(prev_elt_addr, elt_addr) != 0 {
				events <- &ChangedMember{op.name, key_coll, key_addr, prev_elt_addr, elt_addr}
			}

		case UNSET:
			var (
				key_coll []byte
				key_addr cas.Addr
				elt_addr cas.Addr
				err      error
			)

			key_coll = cas.Collate(change.Key)

			key_addr, elt_addr, err = table.Del(key_coll)
			if err != nil {
				panic("runtime: " + err.Error())
			}

			if key_addr != nil || elt_addr != nil {
				events <- &ChangedMember{op.name, key_coll, key_addr, elt_addr, nil}
			}

		}
	}

	tab_addr_a, tab_addr_b := state.CommitTable(op.name, table)
	events <- &ConsistentTable{op.name, tab_addr_a, tab_addr_b}
}

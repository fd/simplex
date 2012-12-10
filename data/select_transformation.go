package data

type SelectFunc func(Context, Value) bool

func Select(f SelectFunc) View {
	return current_engine.ScopedView().Select(f)
}

func (v View) Select(f SelectFunc) View {
	return v.add_transformation(&select_transformation{
		f: f,
	})
}

type select_transformation struct {
	f SelectFunc
	s *select_state
}

type select_state struct {
	Ids []string
}

func (t *select_transformation) Transform(txn transaction) {
	selected := make(map[string]bool, len(t.s.Ids))
	upstream := txn.upstream_states[0]

	for _, id := range t.s.Ids {
		selected[id] = true
	}

	for _, id := range txn.added {
		val := upstream.Get(id)

		if t.f(Context{Id: id}, val) {
			selected[id] = true
		}
	}

	for _, id := range txn.updated {
		val := upstream.Get(id)
		selected[id] = t.f(Context{Id: id}, val)
	}

	for _, id := range txn.removed {
		selected[id] = false
	}

	ids := make([]string, 0, len(selected))

	for _, id := range upstream.Ids() {
		if selected[id] {
			ids = append(ids, id)
		}
	}

	t.s.Ids = ids
}

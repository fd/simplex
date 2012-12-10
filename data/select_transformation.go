package data

type SelectFunc func(Document) bool

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
	Ids []int
}

func (t *select_transformation) Transform(prev State, txn transaction) {
	selected := make(map[int]bool, len(t.s.Ids))

	for _, id := range t.s.Ids {
		selected[id] = true
	}

	for _, id := range txn.added {
		val := prev.Get(id)

		if t.f(val) {
			selected[id] = true
		}
	}

	for _, id := range txn.updated {
		val := prev.Get(id)
		selected[id] = t.f(val)
	}

	for _, id := range txn.removed {
		selected[id] = false
	}

	ids := make([]int, 0, len(selected))

	for _, id := range prev.Ids() {
		if selected[id] {
			ids = append(ids, id)
		}
	}

	t.s.Ids = ids
}

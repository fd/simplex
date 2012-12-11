package data

type SelectFunc func(Context, Value) bool

func Select(f SelectFunc) View {
	return current_engine.ScopedView().Select(f)
}

func (v View) Select(f SelectFunc) View {
	return v.push(&select_transformation{
		id: v.new_id(),
		b:  v.current,
		f:  f,
	})
}

type select_transformation struct {
	id          string
	b           transformation
	f           SelectFunc
	SelectedIds []string
}

func (t *select_transformation) Id() string {
	return t.id
}

func (t *select_transformation) Transform(txn transaction) {
	selected := make(map[string]bool, len(t.SelectedIds))

	for _, id := range t.SelectedIds {
		selected[id] = true
	}

	for _, id := range txn.added {
		val := t.b.Get(id)

		if t.f(Context{Id: id}, val) {
			selected[id] = true
		}
	}

	for _, id := range txn.updated {
		val := t.b.Get(id)
		selected[id] = t.f(Context{Id: id}, val)
	}

	for _, id := range txn.removed {
		selected[id] = false
	}

	ids := make([]string, 0, len(selected))

	for _, id := range t.b.Ids() {
		if selected[id] {
			ids = append(ids, id)
		}
	}

	t.SelectedIds = ids
}

func (t *select_transformation) Restore(txn transaction) {
	txn.state.Restore(t.id, t)
}

func (t *select_transformation) Save(txn transaction) {
	txn.state.Save(t.id, t)
}

func (t *select_transformation) Ids() []string {
	return t.SelectedIds
}

func (t *select_transformation) Get(id string) Value {
	return t.b.Get(id)
}

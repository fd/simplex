package data

func Window(offset, limit int) View {
	return current_engine.ScopedView().Window(offset, limit)
}

func (v View) Window(offset, limit int) View {
	return v.push(&window_transformation{
		id:     v.new_id(),
		b:      v.current,
		offset: offset,
		limit:  limit,
	})
}

func (v View) Offset(n int) View {
	return v.Window(n, 0)
}

func (v View) Limit(n int) View {
	return v.Window(0, n)
}

type window_transformation struct {
	id          string
	b           transformation
	offset      int
	limit       int
	SelectedIds []string
}

func (t *window_transformation) Id() string {
	return t.id
}

func (t *window_transformation) Transform(txn transaction) {
	ids := t.b.Ids()

	if t.limit == 0 {
		ids = ids[t.offset:]
	} else {
		ids = ids[t.offset:t.limit]
	}

	t.SelectedIds = ids
}
func (t *window_transformation) Restore(txn transaction) {
	txn.state.Restore(t.id, t)
}

func (t *window_transformation) Save(txn transaction) {
	txn.state.Save(t.id, t)
}

func (t *window_transformation) Ids() []string {
	return t.SelectedIds
}

func (t *window_transformation) Get(id string) Value {
	return t.Get(id)
}

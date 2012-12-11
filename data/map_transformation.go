package data

type MapFunc func(Context, Value) Value

func Map(f MapFunc) View {
	return current_engine.ScopedView().Map(f)
}

func (v View) Map(f MapFunc) View {
	return v.push(&map_transformation{
		transformation_info: &transformation_info{},
		base:                v.current.Info().Id,
		f:                   f,
	})
}

type map_transformation struct {
	*transformation_info
	base string
	f    MapFunc
}

type map_transformation_state struct {
	base   transformation_state
	Values map[string]Value
}

func (t *map_transformation) Transform(txn transaction) {
	state := &map_transformation_state{
		base: txn.GetStore(t.base),
	}
	txn.Restore(t.Id, &state)

	for _, id := range state.base.Added() {
		val := state.base.Get(id)
		state.Values[id] = t.f(Context{Id: id}, val)
	}

	for _, id := range state.base.Updated() {
		val := state.base.Get(id)
		state.Values[id] = t.f(Context{Id: id}, val)
	}

	for _, id := range state.base.Removed() {
		delete(state.Values, id)
	}

	txn.Save(t.Id, &state)
}

func (t *map_transformation_state) Ids() []string {
	return t.base.Ids()
}

func (t *map_transformation_state) Get(id string) Value {
	return t.Values[id]
}

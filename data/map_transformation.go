package data

type MapFunc func(Context, Value) Value

func Map(f MapFunc) View {
	return current_engine.ScopedView().Map(f)
}

func (v View) Map(f MapFunc) View {
	return v.add_transformation(&map_transformation{
		f: f,
	})
}

type map_transformation struct {
	f MapFunc
	s *map_state
}

type map_state struct {
	Values map[string]Value
}

func (t *map_transformation) Transform(txn transaction) {
	upstream := txn.upstream_states[0]

	for _, id := range txn.added {
		val := upstream.Get(id)
		t.s.Values[id] = t.f(Context{Id: id}, val)
	}

	for _, id := range txn.updated {
		val := upstream.Get(id)
		t.s.Values[id] = t.f(Context{Id: id}, val)
	}

	for _, id := range txn.removed {
		delete(t.s.Values, id)
	}

}

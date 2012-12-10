package data

type MapFunc func(Context, Value) Value

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
	Values map[int]Value
}

func (t *map_transformation) Transform(prev State, txn transaction) {

	for _, id := range txn.added {
		val := prev.Get(id)
		t.s.Values[id] = t.f(val)
	}

	for _, id := range txn.updated {
		val := prev.Get(id)
		t.s.Values[id] = t.f(val)
	}

	for _, id := range txn.removed {
		delete(t.s.Values, id)
	}

}

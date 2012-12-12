package data

type MapFunc func(Context, Value) Value

func Map(f MapFunc) View {
	return current_engine.ScopedView().Map(f)
}

func (v View) Map(f MapFunc) View {
	return v.push(&map_transformation{
		id:       v.new_id() + ":Map",
		upstream: v.current,
		f:        f,
	})
}

type map_transformation struct {
	id         string
	upstream   transformation
	downstream []transformation
	f          MapFunc
}

func (t *map_transformation) Id() string {
	return t.id
}

func (t *map_transformation) Chain() []transformation {
	if t.upstream == nil {
		return []transformation{t}
	}
	return append(t.upstream.Chain(), t)
}

func (t *map_transformation) Dependencies() []transformation {
	if t.upstream == nil {
		return []transformation{}
	}
	return append(t.upstream.Dependencies(), t.upstream)
}

func (t *map_transformation) PushDownstream(d transformation) {
	t.downstream = append(t.downstream, d)
}

func (t *map_transformation) Transform(upstream upstream_state, txn *transaction) {
	var (
		state = upstream.NewState(t.id)
		info  = &map_transformation_state{}
	)

	info.upstream = upstream
	txn.Restore(state, &info)
	state.Info = info

	if info.Values == nil {
		info.Values = make(map[string]Value)
	}

	{
		for _, id := range upstream.Added() {
			val := upstream.Get(id)
			info.Values[id] = t.f(Context{Id: id}, val)
		}

		for _, id := range upstream.Changed() {
			val := upstream.Get(id)
			info.Values[id] = t.f(Context{Id: id}, val)
		}

		for _, id := range upstream.Removed() {
			delete(info.Values, id)
		}
	}

	txn.Save(state)
	txn.Propagate(t.downstream, state)
}

type map_transformation_state struct {
	upstream upstream_state
	Values   map[string]Value
}

func (t *map_transformation_state) Ids() []string {
	return t.upstream.Ids()
}

func (t *map_transformation_state) Get(id string) Value {
	return t.Values[id]
}

package data

func Window(offset, limit int) View {
	return current_engine.ScopedView().Window(offset, limit)
}

func (v View) Window(offset, limit int) View {
	return v.push(&window_transformation{
		id:       v.new_id() + ":Window",
		upstream: v.current,
		offset:   offset,
		limit:    limit,
	})
}

func (v View) Offset(n int) View {
	return v.Window(n, 0)
}

func (v View) Limit(n int) View {
	return v.Window(0, n)
}

type window_transformation struct {
	id         string
	upstream   transformation
	downstream []transformation
	offset     int
	limit      int
}

func (t *window_transformation) Id() string {
	return t.id
}

func (t *window_transformation) Chain() []transformation {
	if t.upstream == nil {
		return []transformation{t}
	}
	return append(t.upstream.Chain(), t)
}

func (t *window_transformation) Dependencies() []transformation {
	if t.upstream == nil {
		return []transformation{}
	}
	return append(t.upstream.Dependencies(), t.upstream)
}

func (t *window_transformation) PushDownstream(d transformation) {
	t.downstream = append(t.downstream, d)
}

func (t *window_transformation) Transform(upstream upstream_state, txn *transaction) {
	var (
		state = upstream.NewState(t.id)
		info  = &window_transformation_state{}
	)

	info.upstream = upstream
	txn.Restore(state, &info)
	state.Info = info

	{
		ids := upstream.Ids()

		if t.limit == 0 {
			ids = ids[t.offset:]
		} else {
			ids = ids[t.offset:t.limit]
		}

		info.SelectedIds = ids
	}

	txn.Save(state)
	txn.Propagate(t.downstream, state)
}

type window_transformation_state struct {
	upstream    upstream_state
	SelectedIds []string
}

func (s *window_transformation_state) Ids() []string {
	return s.SelectedIds
}

func (s *window_transformation_state) Get(id string) Value {
	return s.upstream.Get(id)
}

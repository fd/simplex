package data

type WhereFunc func(Context, Value) bool

func Where(f WhereFunc) View {
	return current_engine.ScopedView().Where(f)
}

func (v View) Where(f WhereFunc) View {
	return v.push(&where_transformation{
		id:       v.new_id() + ":Where",
		upstream: v.current,
		f:        f,
	})
}

type where_transformation struct {
	id         string
	upstream   transformation
	downstream []transformation
	f          WhereFunc
}

func (t *where_transformation) Id() string {
	return t.id
}

func (t *where_transformation) Chain() []transformation {
	if t.upstream == nil {
		return []transformation{t}
	}
	return append(t.upstream.Chain(), t)
}

func (t *where_transformation) Dependencies() []transformation {
	if t.upstream == nil {
		return []transformation{}
	}
	return append(t.upstream.Dependencies(), t.upstream)
}

func (t *where_transformation) PushDownstream(d transformation) {
	t.downstream = append(t.downstream, d)
}

func (t *where_transformation) Transform(upstream upstream_state, txn *transaction) {
	var (
		state = upstream.NewState(t.id)
		info  = &where_transformation_state{}
	)

	info.upstream = upstream
	txn.Restore(state, &info)
	state.Info = info

	{
		selected := make(map[string]bool, len(info.SelectedIds))

		for _, id := range info.SelectedIds {
			selected[id] = true
		}

		for _, id := range upstream.Added() {
			val := upstream.Get(id)

			if t.f(Context{Id: id}, val) {
				selected[id] = true
				state.added = append(state.added, id)
			}
		}

		for _, id := range upstream.Changed() {
			val := upstream.Get(id)
			sel := t.f(Context{Id: id}, val)
			if selected[id] == true && sel == true {
				state.changed = append(state.changed, id)
			} else if selected[id] == true && sel == false {
				state.removed = append(state.removed, id)
			} else if selected[id] == false && sel == true {
				state.added = append(state.added, id)
			}
			selected[id] = sel
		}

		for _, id := range upstream.Removed() {
			if selected[id] {
				selected[id] = false
				state.removed = append(state.removed, id)
			}
		}

		ids := make([]string, 0, len(selected))

		for _, id := range upstream.Ids() {
			if selected[id] {
				ids = append(ids, id)
			}
		}

		info.SelectedIds = ids
	}

	txn.Save(state)
	txn.Propagate(t.downstream, state)
}

type where_transformation_state struct {
	upstream    upstream_state
	SelectedIds []string
}

func (s *where_transformation_state) Ids() []string {
	return s.SelectedIds
}

func (s *where_transformation_state) Get(id string) Value {
	return s.upstream.Get(id)
}

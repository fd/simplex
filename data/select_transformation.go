package data

type SelectFunc func(Context, Value) bool

func Select(f SelectFunc) View {
	return current_engine.ScopedView().Select(f)
}

func (v View) Select(f SelectFunc) View {
	return v.push(&select_transformation{
		id:       v.new_id() + ":Select",
		upstream: v.current,
		f:        f,
	})
}

type select_transformation struct {
	id         string
	upstream   transformation
	downstream []transformation
	f          SelectFunc
}

func (t *select_transformation) Id() string {
	return t.id
}

func (t *select_transformation) Chain() []transformation {
	if t.upstream == nil {
		return []transformation{t}
	}
	return append(t.upstream.Chain(), t)
}

func (t *select_transformation) Dependencies() []transformation {
	if t.upstream == nil {
		return []transformation{}
	}
	return append(t.upstream.Dependencies(), t.upstream)
}

func (t *select_transformation) PushDownstream(d transformation) {
	t.downstream = append(t.downstream, d)
}

func (t *select_transformation) Transform(upstream upstream_state, txn *transaction) {
	var (
		state = upstream.NewState(t.id)
		info  = &select_transformation_state{}
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

type select_transformation_state struct {
	upstream    upstream_state
	SelectedIds []string
}

func (s *select_transformation_state) Ids() []string {
	return s.SelectedIds
}

func (s *select_transformation_state) Get(id string) Value {
	return s.upstream.Get(id)
}

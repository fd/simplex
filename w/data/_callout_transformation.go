package data

type CalloutFunc func(Context, Value) []string

func Callout(members View, f CalloutFunc) View {
	return current_engine.ScopedView().Callout(members, f)
}

func (v View) Callout(members View, f CalloutFunc) View {
	return v.push(&callout_transformation{
		transformation_info: &transformation_info{},
		base:                v.current.Info().Id,
		members:             members.current.Info().Id,
		f:                   f,
	})
}

type callout_transformation struct {
	*transformation_info
	base    string
	members string
	f       CalloutFunc
}

type callout_transformation_state struct {
	base       transformation_state
	members    transformation_state
	CalloutIds map[string][]string
}

func (t *callout_transformation) Id() string {
	return t.id
}

func (t *callout_transformation) Transform(txn transaction) {
	state := &callout_transformation_state{
		base:    txn.GetStore(t.base),
		members: txn.GetStore(t.members),
	}
	txn.Restore(t.Id, &state)

	for _, id := range txn.added {
		val := t.base.Get(id)
		t.CalloutIds[id] = t.f(Context{Id: id}, val)
	}

	for _, id := range txn.updated {
		val := t.base.Get(id)
		t.CalloutIds[id] = t.f(Context{Id: id}, val)
	}

	for _, id := range txn.removed {
		delete(t.CalloutIds, id)
	}

	txn.Save(t.Id, &state)
}

func (t *callout_transformation) Restore(txn transaction) {
	txn.state.Restore(t.id, t)
}

func (t *callout_transformation) Save(txn transaction) {
	txn.state.Save(t.id, t)
}

type callout struct {
	t          *callout_transformation
	CalloutIds []string
}

func (t *callout) Ids() []string {
	return t.CalloutIds
}

func (t *callout) Get(id string) Value {
	return t.t.members.Get(id)
}

package data

type Changes struct {
	Create  map[string]Value
	Update  map[string]Value
	Destroy []string
	engine  *Engine
}

func (c Changes) Added() []string {
	ids := make([]string, 0, len(c.Create))
	for id := range c.Create {
		ids = append(ids, id)
	}
	return ids
}

func (c Changes) Changed() []string {
	ids := make([]string, 0, len(c.Update))
	for id := range c.Update {
		ids = append(ids, id)
	}
	return ids
}

func (c Changes) Removed() []string {
	return c.Destroy
}

func (c Changes) Ids() []string {
	return c.engine.SourceTable.Ids()
}

func (c Changes) Get(id string) Value {
	return c.engine.SourceTable.Get(id)
}

func (c Changes) NewState(segment ...string) *state {
	return &state{
		Id: append([]string{"$root"}, segment...),
	}
}

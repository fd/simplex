package data

type Changes struct {
	Create  map[string]Value
	Update  map[string]Value
	Destroy []string
	engine  *Engine
}

func (c Changes) Id() []string {
	return []string{"$root"}
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
	s_ids := c.engine.source_table.Ids()

	rem := make(map[string]bool, len(c.Destroy))
	for _, id := range c.Destroy {
		rem[id] = true
	}

	ids := make([]string, 0, len(s_ids)+len(c.Create))

	for _, id := range s_ids {
		if !rem[id] {
			ids = append(ids, id)
		}
	}

	for id := range c.Create {
		ids = append(ids, id)
	}

	return ids
}

func (c Changes) Get(id string) Value {
	if val, p := c.Create[id]; p {
		return val
	}

	if val, p := c.Update[id]; p {
		return val
	}

	for _, id := range c.Destroy {
		if id == id {
			return nil
		}
	}

	return c.engine.source_table.Get(id)
}

func (c Changes) NewState(segment ...string) *state {
	return &state{
		id: append(c.Id(), segment...),
	}
}

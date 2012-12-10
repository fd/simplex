package data

func Window(offset, limit int) View {
	return current_engine.ScopedView().Window(offset, limit)
}

func (v View) Window(offset, limit int) View {
	return v.add_transformation(&window_transformation{
		offset: offset,
		limit:  limit,
	})
}

func (v View) Offset(n int) View {
	return v.Window(n, 0)
}

func (v View) Limit(n int) View {
	return v.Window(0, n)
}

type window_transformation struct {
	offset int
	limit  int
	s      *window_state
}

type window_state struct {
	Ids []string
}

func (t *window_transformation) Transform(txn transaction) {
	upstream := txn.upstream_states[0]

	ids := upstream.Ids()

	if t.limit == 0 {
		ids = ids[t.offset:]
	} else {
		ids = ids[t.offset:t.limit]
	}

	t.s.Ids = ids
}

package data

import (
	"sort"
)

type SortFunc func(Context, Value) Value

func Sort(f SortFunc) View {
	return current_engine.ScopedView().Sort(f)
}

func (v View) Sort(f SortFunc) View {
	return v.add_transformation(&sort_transformation{
		f: f,
	})
}

type sort_transformation struct {
	f SortFunc
	s *sort_state
}

type sort_state struct {
	Ids  []string
	Keys map[string]Value
}

func (t *sort_transformation) Transform(txn transaction) {
	upstream := txn.upstream_states[0]

	for _, id := range txn.added {
		val := upstream.Get(id)
		key := t.f(Context{Id: id}, val)
		t.s.Keys[id] = key
	}

	for _, id := range txn.updated {
		val := upstream.Get(id)
		key := t.f(Context{Id: id}, val)
		t.s.Keys[id] = key
	}

	for _, id := range txn.removed {
		delete(t.s.Keys, id)
	}

	ids := make([]string, 0, len(t.s.Keys))
	for id := range t.s.Keys {
		ids = append(ids, id)
	}
	t.s.Ids = ids

	sort.Sort(t.s)
}

func (s *sort_state) Len() int {
	return len(s.Ids)
}

func (s *sort_state) Less(i, j int) bool {
	m, n := s.Ids[i], s.Ids[j]
	x, y := s.Keys[m], s.Keys[n]
	return Compair(x, y) == -1
}

func (s *sort_state) Swap(i, j int) {
	s.Ids[j], s.Ids[i] = s.Ids[i], s.Ids[j]
}

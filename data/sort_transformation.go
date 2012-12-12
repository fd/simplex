package data

import (
	"sort"
)

type SortFunc func(Context, Value) Value

func Sort(f SortFunc) View {
	return current_engine.ScopedView().Sort(f)
}

func (v View) Sort(f SortFunc) View {
	return v.push(&sort_transformation{
		id:       v.new_id() + ":Sort",
		upstream: v.current,
		f:        f,
	})
}

type sort_transformation struct {
	id         string
	upstream   transformation
	downstream []transformation
	f          SortFunc
}

func (t *sort_transformation) Id() string {
	return t.id
}

func (t *sort_transformation) Chain() []transformation {
	if t.upstream == nil {
		return []transformation{t}
	}
	return append(t.upstream.Chain(), t)
}

func (t *sort_transformation) Dependencies() []transformation {
	if t.upstream == nil {
		return []transformation{}
	}
	return append(t.upstream.Dependencies(), t.upstream)
}

func (t *sort_transformation) PushDownstream(d transformation) {
	t.downstream = append(t.downstream, d)
}

func (t *sort_transformation) Transform(upstream upstream_state, txn *transaction) {
	var (
		state = upstream.NewState(t.id)
		info  = &sort_transformation_state{}
	)

	info.upstream = upstream
	txn.Restore(state, &info)
	state.Info = info

	if info.Keys == nil {
		info.Keys = make(map[string]string)
	}

	{
		state.added = upstream.Added()
		state.changed = upstream.Changed()
		state.removed = upstream.Removed()

		for _, id := range upstream.Added() {
			val := upstream.Get(id)
			key_val := t.f(Context{Id: id}, val)
			key_str := CompairString(key_val)
			info.Keys[id] = key_str
		}

		for _, id := range upstream.Changed() {
			val := upstream.Get(id)
			key_val := t.f(Context{Id: id}, val)
			key_str := CompairString(key_val)
			info.Keys[id] = key_str
		}

		for _, id := range upstream.Removed() {
			delete(info.Keys, id)
		}

		ids := make([]string, 0, len(info.Keys))
		for id := range info.Keys {
			ids = append(ids, id)
		}
		info.SortedIds = ids
		sort.Sort(info)
	}

	txn.Save(state)
	txn.Propagate(t.downstream, state)
}

type sort_transformation_state struct {
	upstream  upstream_state
	SortedIds []string
	Keys      map[string]string
}

func (s *sort_transformation_state) Len() int {
	return len(s.SortedIds)
}

func (s *sort_transformation_state) Less(i, j int) bool {
	m, n := s.SortedIds[i], s.SortedIds[j]
	x, y := s.Keys[m], s.Keys[n]
	if x == y {
		return m < n
	}
	return x < y
}

func (s *sort_transformation_state) Swap(i, j int) {
	s.SortedIds[j], s.SortedIds[i] = s.SortedIds[i], s.SortedIds[j]
}

func (s *sort_transformation_state) Ids() []string {
	return s.SortedIds
}

func (s *sort_transformation_state) Get(id string) Value {
	return s.upstream.Get(id)
}

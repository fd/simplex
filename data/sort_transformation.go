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
		id: v.new_id(),
		b:  v.current,
		f:  f,
	})
}

type sort_transformation struct {
	id        string
	b         transformation
	f         SortFunc
	SortedIds []string
	Keys      map[string]Value
}

func (t *sort_transformation) Id() string {
	return t.id
}

func (t *sort_transformation) Transform(txn transaction) {

	for _, id := range txn.added {
		val := t.b.Get(id)
		key := t.f(Context{Id: id}, val)
		t.Keys[id] = key
	}

	for _, id := range txn.updated {
		val := t.b.Get(id)
		key := t.f(Context{Id: id}, val)
		t.Keys[id] = key
	}

	for _, id := range txn.removed {
		delete(t.Keys, id)
	}

	ids := make([]string, 0, len(t.Keys))
	for id := range t.Keys {
		ids = append(ids, id)
	}
	t.SortedIds = ids

	sort.Sort(t)
}

func (s *sort_transformation) Len() int {
	return len(s.SortedIds)
}

func (s *sort_transformation) Less(i, j int) bool {
	m, n := s.SortedIds[i], s.SortedIds[j]
	x, y := s.Keys[m], s.Keys[n]
	return Compair(x, y) == -1
}

func (s *sort_transformation) Swap(i, j int) {
	s.SortedIds[j], s.SortedIds[i] = s.SortedIds[i], s.SortedIds[j]
}

func (t *sort_transformation) Restore(txn transaction) {
	txn.state.Restore(t.id, t)
}

func (t *sort_transformation) Save(txn transaction) {
	txn.state.Save(t.id, t)
}

func (t *sort_transformation) Ids() []string {
	return t.SortedIds
}

func (t *sort_transformation) Get(id string) Value {
	return t.b.Get(id)
}

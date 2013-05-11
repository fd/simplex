package static

import (
	"sort"
)

// Sort a collection. The function f must return true if a is less than b (a < b).
func (in *C) Sort(f func(a, b interface{}) bool) *C {
	return in.Transform(in.elem_type, func(i_elems []interface{}) ([]interface{}, error) {
		var (
			o_elems = make([]interface{}, len(i_elems))
		)

		copy(o_elems, i_elems)
		sort.Sort(&collection_sorter{o_elems, f})

		return o_elems, nil
	})
}

type collection_sorter struct {
	slice []interface{}
	less  func(a, b interface{}) bool
}

func (s *collection_sorter) Len() int {
	return len(s.slice)
}

func (s *collection_sorter) Less(i, j int) bool {
	return s.less(s.slice[i], s.slice[j])
}

func (s *collection_sorter) Swap(i, j int) {
	s.slice[i], s.slice[j] = s.slice[j], s.slice[i]
}

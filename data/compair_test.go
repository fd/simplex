package data

import (
	"sort"
	"testing"
)

func Test_compair(t *testing.T) {
	v := []Value{
		// special values sort before all other types
		nil,
		false,
		true,

		// then numbers
		1,
		2,
		3.0,
		4,

		// then text, case sensitive
		"A",
		"B",
		"a",
		"aa",
		"b",
		"ba",
		"bb",

		// then arrays. compared element by element until different.
		// Longer arrays sort after their prefixes
		[]Value{"a"},
		[]Value{"b"},
		[]Value{"b", "c"},
		[]Value{"b", "c", "a"},
		[]Value{"b", "d"},
		[]Value{"b", "d", "e"},

		// then object, compares each key value in the list until different.
		// larger objects sort after their subset objects.
		map[string]Value{"a": 1},
		map[string]Value{"a": 2},
		map[string]Value{"b": 1},
		map[string]Value{"b": 2},
		map[string]Value{"b": 2, "c": 2},
	}

	r := make([]Value, len(v))
	for i, e := range v {
		r[len(v)-1-i] = e
	}

	sort.Sort(sort_slice(r))

	for i, a := range r {
		e := v[i]
		if CompairString(a) != CompairString(e) {
			t.Errorf("Expected %+v but was %+v", e, a)
		}
	}
}

type sort_slice []Value

func (s sort_slice) Len() int {
	return len(s)
}

func (s sort_slice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sort_slice) Less(i, j int) bool {
	return CompairString(s[i]) < CompairString(s[j])
}

package ident

import (
	"sort"
	"testing"
)

func Test_compair(t *testing.T) {
	v := []interface{}{
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
		[]interface{}{"a"},
		[]interface{}{"b"},
		[]interface{}{"b", "c"},
		[]interface{}{"b", "c", "a"},
		[]interface{}{"b", "d"},
		[]interface{}{"b", "d", "e"},

		// then object, compares each key value in the list until different.
		// larger objects sort after their subset objects.
		map[string]interface{}{"a": 1},
		map[string]interface{}{"a": 2},
		map[string]interface{}{"b": 1},
		map[string]interface{}{"b": 2},
		map[string]interface{}{"b": 2, "c": 2},
	}

	r := make([]interface{}, len(v))
	for i, e := range v {
		r[len(v)-1-i] = e
	}

	sort.Sort(sort_slice(r))

	for i, a := range r {
		e := v[i]
		if string(CompairBytes(a)) != string(CompairBytes(e)) {
			t.Errorf("Expected %+v but was %+v", e, a)
		}
	}
}

type sort_slice []interface{}

func (s sort_slice) Len() int {
	return len(s)
}

func (s sort_slice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sort_slice) Less(i, j int) bool {
	return string(CompairBytes(s[i])) < string(CompairBytes(s[j]))
}

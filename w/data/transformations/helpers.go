package transformations

import (
	"github.com/fd/simplex/w/data/value"
)

type WhereFunc func(key, val value.Any) bool

func Where(f WhereFunc) {
	mr := &MapReduce{
		Map: func(e Emiter, key, val value.Any) {
			if f(key, val) {
				e.Emit(key, val)
			}
		},
		Reduce: IdentifyReduceFunc,
	}
}

type MapFunc func(key, val value.Any) (key, val value.Any)

func Map(f MapFunc) {
	mr := &MapReduce{
		Map: func(e Emiter, key, val value.Any) {
			e.Emit(f(key, val))
		},
		Reduce: IdentifyReduceFunc,
	}
}

type SortFunc func(key, val value.Any) (sort_key value.Any)

func Sort(f SortFunc) {
	mr := &MapReduce{
		Map: func(e Emiter, key, val value.Any) {
			e.Emit([]value.Any{f(key, val), key}, val)
		},
		Reduce: IdentifyReduceFunc,
	}
}

type GroupFunc func(key, val value.Any) (group_key value.Any)

func Group(f SortFunc) {
	mr := &MapReduce{
		Map: func(e Emiter, key, val value.Any) {
			e.Emit(f(key, val), &group_entry{key, val})
		},
		Reduce: func(e Emiter, key value.Any, vals []value.Any) {
			// this needs the reduce trie
		},
	}

	co := Callout{
		F: func(key, val value.Any) (context value.Any, members *trie.T) {
			return CompairBytes(key), val
		},
	}
}

type group_entry struct {
	Key value.Any
	Val value.Any
}

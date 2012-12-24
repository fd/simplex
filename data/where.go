package data

type WhereFunc func(key, val Value) bool

func Where(f WhereFunc) {
	mr := &MapReduce{
		Map: func(e Emiter, key, val Value) {
			if f(key, val) {
				e.Emit(key, val)
			}
		},
		Reduce: IdentifyReduceFunc,
	}
}

type MapFunc func(key, val Value) (key, val Value)

func Map(f MapFunc) {
	mr := &MapReduce{
		Map: func(e Emiter, key, val Value) {
			e.Emit(f(key, val))
		},
		Reduce: IdentifyReduceFunc,
	}
}

type SortFunc func(key, val Value) (sort_key Value)

func Sort(f SortFunc) {
	mr := &MapReduce{
		Map: func(e Emiter, key, val Value) {
			e.Emit([]Value{f(key, val), key}, val)
		},
		Reduce: IdentifyReduceFunc,
	}
}

type GroupFunc func(key, val Value) (group_key Value)

func Group(f SortFunc) {
	mr := &MapReduce{
		Map: func(e Emiter, key, val Value) {
			e.Emit(f(key, val), &group_entry{key, val})
		},
		Reduce: func(e Emiter, key Value, vals []Value) {
			// this needs the reduce trie
		},
	}

	co := Callout{
		F: func(key, val Value) (context Value, members *trie.T) {
			return CompairBytes(key), val
		},
	}
}

type group_entry struct {
	Key Value
	Val Value
}

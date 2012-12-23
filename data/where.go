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

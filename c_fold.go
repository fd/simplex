package static

func (in *C) Fold(init interface{}, f func(acc, v interface{}) interface{}) *Promise {
	var (
		p = &Promise{tx: in.tx}
	)

	p.Do(func() (interface{}, error) {
		if err := in.t.Wait(); err != nil {
			return nil, err
		}

		var (
			acc = init
		)

		for _, elem := range in.elems {
			acc = f(acc, elem)
		}

		return acc, nil
	})

	return p
}

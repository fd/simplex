package static

func (in *C) PromiseAt(idx int) *Promise {
	p := &Promise{tx: in.tx}

	p.Do(func() (interface{}, error) {
		if err := in.t.Wait(); err != nil {
			return nil, err
		}

		if idx < len(in.elems) {
			return in.elems[idx], nil
		}

		return nil, nil
	})

	return p
}

func (in *C) At(idx int) (interface{}, error) {
	return in.PromiseAt(idx).Wait()
}

func (in *C) MustAt(idx int) interface{} {
	v, err := in.At(idx)
	if err != nil {
		panic(err)
	}
	return v
}

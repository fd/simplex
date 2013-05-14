package static

import (
	"simplex.sh/future"
)

func (in *C) PromiseFold(init interface{}, f func(acc, v interface{}) interface{}) future.P {
	var (
		p = &future.Promise{}
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

func (in *C) Fold(init interface{}, f func(acc, v interface{}) interface{}) (interface{}, error) {
	return in.PromiseFold(init, f).Wait()
}

func (in *C) MustFold(init interface{}, f func(acc, v interface{}) interface{}) interface{} {
	v, err := in.Fold(init, f)
	if err != nil {
		panic(err)
	}
	return v
}

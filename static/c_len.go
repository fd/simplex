package static

import (
	"simplex.sh/future"
)

func (in *C) PromiseLen() future.P {
	p := &future.Promise{}

	p.Do(func() (interface{}, error) {
		if err := in.t.Wait(); err != nil {
			return nil, err
		}

		return len(in.elems), nil
	})

	return p
}

func (in *C) Len() (int, error) {
	v, err := in.PromiseLen().Wait()
	if err != nil {
		return 0, err
	}
	return v.(int), nil
}

func (in *C) MustLen() int {
	v, err := in.Len()
	if err != nil {
		panic(err)
	}
	return v
}

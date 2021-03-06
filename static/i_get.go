package static

import (
	"simplex.sh/future"
)

func (in *I) PromiseGet(key string) future.P {
	p := &future.Promise{}

	p.Do(func() (interface{}, error) {
		if err := in.t.Wait(); err != nil {
			return nil, err
		}

		return in.elems[key], nil
	})

	return p
}

func (in *I) Get(key string) (interface{}, error) {
	return in.PromiseGet(key).Wait()
}

func (in *I) MustGet(key string) interface{} {
	v, err := in.Get(key)
	if err != nil {
		panic(err)
	}
	return v
}

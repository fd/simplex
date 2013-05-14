package static

import (
	"simplex.sh/future"
)

func (in *G) PromiseLen(key string) future.P {
	p := &future.Promise{}

	p.Do(func() (interface{}, error) {
		if err := in.t.Wait(); err != nil {
			return nil, err
		}

		return len(in.elems), nil
	})

	return p
}

func (in *G) Len(key string) (int, error) {
	v, err := in.PromiseLen(key).Wait()
	if err != nil {
		return 0, err
	}
	return v.(int), nil
}

func (in *G) MustLen(key string) int {
	v, err := in.Len(key)
	if err != nil {
		panic(err)
	}
	return v
}

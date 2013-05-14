package static

import (
	"simplex.sh/future"
)

func (in *C) SumInt(f func(v interface{}) int64) future.P {
	return in.PromiseFold(0, func(acc, v interface{}) interface{} { return acc.(int64) + f(v) })
}

func (in *C) SumFloat(f func(v interface{}) float64) future.P {
	return in.PromiseFold(0, func(acc, v interface{}) interface{} { return acc.(float64) + f(v) })
}

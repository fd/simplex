package static

func (in *C) SumInt(f func(v interface{}) int64) *Promise {
	return in.Fold(0, func(acc, v interface{}) interface{} { return acc.(int64) + f(v) })
}

func (in *C) SumFloat(f func(v interface{}) float64) *Promise {
	return in.Fold(0, func(acc, v interface{}) interface{} { return acc.(float64) + f(v) })
}

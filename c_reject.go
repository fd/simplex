package static

func (in *C) Reject(f func(v interface{}) bool) *C {
	return in.Select(func(v interface{}) bool { return !f(v) })
}

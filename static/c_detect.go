package static

func (in *C) PromiseDetect(f func(v interface{}) bool) *Promise {
	return in.Select(f).PromiseAt(0)
}

func (in *C) Detect(f func(v interface{}) bool) (interface{}, error) {
	return in.Select(f).At(0)
}

func (in *C) MustDetect(f func(v interface{}) bool) interface{} {
	return in.Select(f).MustAt(0)
}

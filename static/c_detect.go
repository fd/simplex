package static

import (
	"simplex.sh/future"
)

func (in *C) PromiseDetect(f interface{}) future.P {
	f_type_Detect(in.elem_type, f)
	return in.Select(f).PromiseAt(0)
}

func (in *C) Detect(f interface{}) (interface{}, error) {
	f_type_Detect(in.elem_type, f)
	return in.Select(f).At(0)
}

func (in *C) MustDetect(f interface{}) interface{} {
	f_type_Detect(in.elem_type, f)
	return in.Select(f).MustAt(0)
}

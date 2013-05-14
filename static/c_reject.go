package static

func (in *C) Reject(f interface{}) *C {
	var (
		fv, ft  = f_type_Reject(in.elem_type, f)
		has_err = ft.NumOut() == 2
	)

	return in.select_or_reject(fv, has_err, false)
}

package static

import (
	"reflect"
	"simplex.sh/errors"
)

func (in *C) Select(f interface{}) *C {
	var (
		fv, ft  = f_type_Select(in.elem_type, f)
		has_err = ft.NumOut() == 2
	)

	return in.select_or_reject(fv, has_err, true)
}

func (in *C) select_or_reject(fv reflect.Value, has_err, must_be bool) *C {
	return in.Transform(in.elem_type, func(i_elems []interface{}) ([]interface{}, error) {
		var (
			o_elems  = make([]interface{}, 0, len(i_elems))
			in_args  = make([]reflect.Value, 1)
			out_args []reflect.Value
			err      errors.List
		)

		for _, elem := range i_elems {
			in_args[0] = reflect.ValueOf(elem)
			out_args = fv.Call(in_args)
			if has_err && !out_args[1].IsNil() {
				err.Add(out_args[1].Interface().(error))
				continue
			}
			if out_args[0].Bool() == must_be {
				o_elems = append(o_elems, elem)
			}
		}

		return o_elems, nil
	})
}

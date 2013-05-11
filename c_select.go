package static

import (
	"github.com/fd/static/errors"
	"reflect"
)

func (in *C) Select(f interface{}) *C {
	var (
		fv      = reflect.ValueOf(f)
		ft      = fv.Type()
		it      reflect.Type
		ot      reflect.Type
		has_err bool
	)

	if fv.Kind() != reflect.Func {
		panic("Select(f) expects f to  be a function")
	}

	if ft.NumIn() != 1 {
		panic("Select(f): f must take one argument")
	}

	if ft.NumOut() != 1 && ft.NumOut() != 2 {
		panic("Select(f): f must take one or two argument")
	}

	if ft.NumOut() == 2 {
		if ft.Out(1) != reflect.TypeOf((*error)(nil)).Elem() {
			panic("Select(f): the second return value of f must be of type error")
		}
		has_err = true
	}

	it = ft.In(0)
	ot = ft.Out(0)

	if it.Kind() == reflect.Interface && !in.elem_type.Implements(it) {
		panic("Select(f): the input type of f is not compatible with in")
	} else if it != in.elem_type {
		panic("Select(f): the input type of f is not compatible with in")
	}

	if ot.Kind() != reflect.Bool {
		panic("Select(f): the output type of f must be bool")
	}

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
			if out_args[0].Bool() {
				o_elems = append(o_elems, elem)
			}
		}

		return o_elems, nil
	})
}

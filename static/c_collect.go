package static

import (
	"reflect"
	"runtime"
	"simplex.sh/errors"
	"sync"
)

/*

Collect transforms all the members of in according to f.
f must be a function which takes one argument (Ti or interface{})
and returns one value (To) and an optional error

*/
func (in *C) Collect(f interface{}) *C {
	var (
		fv      = reflect.ValueOf(f)
		ft      = fv.Type()
		it      reflect.Type
		ot      reflect.Type
		has_err bool
	)

	if fv.Kind() != reflect.Func {
		panic("Collect(f) expects f to  be a function")
	}

	if ft.NumIn() != 1 {
		panic("Collect(f): f must take one argument")
	}

	if ft.NumOut() != 1 && ft.NumOut() != 2 {
		panic("Collect(f): f must take one or two argument")
	}

	if ft.NumOut() == 2 {
		if ft.Out(1) != reflect.TypeOf((*error)(nil)).Elem() {
			panic("Collect(f): the second return value of f must be of type error")
		}
		has_err = true
	}

	it = ft.In(0)
	ot = ft.Out(0)

	if it != in.elem_type && it != reflect.TypeOf((*interface{})(nil)).Elem() {
		panic("Collect(f): the input type of f is not compatible with in")
	}

	return in.Transform(ot, func(i_elems []interface{}) ([]interface{}, error) {
		var (
			o_elems    = make([]interface{}, len(i_elems))
			workers    = runtime.NumCPU() * 2
			slice_size = (len(i_elems) / workers) + 1
			ctx        = &collect_context{f: fv, has_err: has_err}
		)

		for i := 0; i < workers; i++ {
			beg := i * slice_size
			end := beg + slice_size

			if end > len(i_elems) {
				end = len(i_elems)
			}

			if beg >= len(i_elems) {
				continue
			}

			ctx.wg.Add(1)
			go go_collect(i_elems[beg:end], o_elems[beg:end], ctx)
		}

		ctx.wg.Wait()
		return o_elems, ctx.err.Normalize()
	})
}

type collect_context struct {
	wg      sync.WaitGroup
	err     errors.List
	f       reflect.Value
	has_err bool
}

func go_collect(in, out []interface{}, ctx *collect_context) {
	defer ctx.wg.Done()

	var (
		in_args  = make([]reflect.Value, 1)
		out_args []reflect.Value
	)

	for i, elem := range in {
		in_args[0] = reflect.ValueOf(elem)
		out_args = ctx.f.Call(in_args)
		if ctx.has_err && !out_args[1].IsNil() {
			ctx.err.Add(out_args[1].Interface().(error))
			continue
		}
		out[i] = out_args[0].Interface()
	}
}

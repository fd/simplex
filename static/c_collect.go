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
	fv, ft := f_type_Collect(in.elem_type, f)

	return in.Transform(ft.Out(0), func(i_elems []interface{}) ([]interface{}, error) {
		var (
			o_elems    = make([]interface{}, len(i_elems))
			workers    = runtime.NumCPU() * 2
			slice_size = len(i_elems)/workers + 1
			ctx        = &collect_context{f: fv, has_err: ft.NumOut() == 2}
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

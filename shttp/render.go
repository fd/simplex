package shttp

import (
	"github.com/fd/static"
	"reflect"
)

var (
	errorType  = reflect.TypeOf((*error)(nil)).Elem()
	writerType = reflect.TypeOf((*Writer)(nil)).Elem()
)

// f: func(m T, w io.Writer, r Router) error
func Render(in *static.C, f interface{}) {
	var (
		fv = reflect.ValueOf(f)
		ft = fv.Type()
		i0 reflect.Type
		i1 reflect.Type
		o0 reflect.Type
	)

	if ft.Kind() != reflect.Func || ft.NumIn() != 2 || ft.NumOut() != 1 {
		panic("Render(f): f must have signature: func(m T, w Writer) error")
	}

	i0 = ft.In(0)
	i1 = ft.In(1)
	o0 = ft.Out(0)

	if !o0.Implements(errorType) || i1 != writerType {
		panic("Render(f): f must have signature: func(m T, w Writer) error")
	}

	if i0.Kind() == reflect.Interface && !in.ElemType().Implements(i0) {
		panic("Render(f): f must have signature: func(m T, w Writer) error")
	} else if i0 != in.ElemType() {
		panic("Render(f): f must have signature: func(m T, w Writer) error")
	}

	router := terminator_for_tx(in.Tx())

	docs := in.Collect(func(v interface{}) (*document, error) {
		var (
			d      *document
			dw     = new_document_writer()
			rw     = Writer(dw)
			args_o []reflect.Value
			args_i = []reflect.Value{
				reflect.ValueOf(v),
				reflect.ValueOf(rw),
			}
		)

		args_o = fv.Call(args_i)

		if !args_o[0].IsNil() {
			return nil, args_o[0].Interface().(error)
		}

		dw.Close()
		d = dw.document

		for _, rule := range dw.route_builder.rules {
			set := router.route_table.path(rule.path)
			err := set.add(&route_rule{
				Host:        rule.host,
				Language:    d.Header.Get("Language"),
				ContentType: d.Header.Get("Content-Type"),
				Status:      d.Status,
				Header:      d.Header,
				Address:     d.Digest,
			})
			if err != nil {
				return nil, err
			}
		}

		return d, nil
	})

	router.mtx.Lock()
	defer router.mtx.Unlock()
	router.collections = append(router.collections, docs)
}

package static

import (
	"reflect"
)

func f_type_Collect(elem_type reflect.Type, f interface{}) (reflect.Value, reflect.Type) {
	var (
		fv = reflect.ValueOf(f)
		ft = fv.Type()
		it reflect.Type
	)

	if ft.Kind() != reflect.Func {
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
	}

	it = ft.In(0)

	if it != elem_type && it != reflect.TypeOf((*interface{})(nil)).Elem() {
		panic("Collect(f): the input type of f is not compatible with in")
	}

	return fv, ft
}

func f_type_Select(elem_type reflect.Type, f interface{}) (reflect.Value, reflect.Type) {
	return f_type_Condition("Select", elem_type, f)
}

func f_type_Reject(elem_type reflect.Type, f interface{}) (reflect.Value, reflect.Type) {
	return f_type_Condition("Reject", elem_type, f)
}

func f_type_Detect(elem_type reflect.Type, f interface{}) (reflect.Value, reflect.Type) {
	return f_type_Condition("Detect", elem_type, f)
}

func f_type_Condition(type_name string, elem_type reflect.Type, f interface{}) (reflect.Value, reflect.Type) {
	var (
		fv = reflect.ValueOf(f)
		ft = fv.Type()
		it reflect.Type
		ot reflect.Type
	)

	if ft.Kind() != reflect.Func {
		panic(type_name + "(f) expects f to  be a function")
	}

	if ft.NumIn() != 1 {
		panic(type_name + "(f): f must take one argument")
	}

	if ft.NumOut() != 1 && ft.NumOut() != 2 {
		panic(type_name + "(f): f must take one or two argument")
	}

	if ft.NumOut() == 2 {
		if ft.Out(1) != reflect.TypeOf((*error)(nil)).Elem() {
			panic(type_name + "(f): the second return value of f must be of type error")
		}
	}

	it = ft.In(0)
	ot = ft.Out(0)

	if it.Kind() == reflect.Interface && !elem_type.Implements(it) {
		panic(type_name + "(f): the input type of f is not compatible with in")
	} else if it != elem_type {
		panic(type_name + "(f): the input type of f is not compatible with in")
	}

	if ot.Kind() != reflect.Bool {
		panic(type_name + "(f): the output type of f must be bool")
	}

	return fv, ft
}

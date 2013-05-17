package collation

import (
	"encoding/base64"
	"encoding/binary"
	"io"
	"reflect"
	"simplex.sh/errors"
)

const (
	low  byte = 0
	high      = 255
)
const (
	nil_type byte = 1 + iota
	bool_false_type
	bool_true_type
	int_type
	uint_type
	float_type
	string_type
	slice_type
	map_type
)

func CollateValue(w io.Writer, v reflect.Value) error {
	switch v.Kind() {

	case reflect.Ptr:
		collate_ptr(w, v)

	case reflect.Bool:
		collate_bool(w, v.Bool())

	case reflect.Int,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		collate_int(w, v.Int())

	case reflect.Uint,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		collate_uint(w, v.Uint())

	case reflect.Float32, reflect.Float64:
		collate_float(w, v.Float())

	case reflect.String:
		collate_string(w, v.String())

	case reflect.Slice, reflect.Array:
		collate_slice(w, v)

	case reflect.Map:
		collate_map(w, v)

	case reflect.Struct:
		collate_struct(w, v)

	default:
		return errors.Fmt("collate: Unable to collate value of type %s", v.Type())

	}

	return nil
}

func collate_ptr(w io.Writer, v reflect.Value) {
	if v.IsNil() {
		w.Write([]byte{nil_type})
	} else {
		CollateValue(w, v.Elem())
	}
}

func collate_bool(w io.Writer, v bool) {
	if v {
		w.Write([]byte{bool_true_type})
	} else {
		w.Write([]byte{bool_false_type})
	}
}

func collate_int(w io.Writer, v int64) {
	w.Write([]byte{int_type})
	binary.Write(w, binary.BigEndian, v)
}

func collate_uint(w io.Writer, v uint64) {
	w.Write([]byte{uint_type})
	binary.Write(w, binary.BigEndian, v)
}

func collate_float(w io.Writer, v float64) {
	w.Write([]byte{uint_type})
	binary.Write(w, binary.BigEndian, v)
}

func collate_string(w io.Writer, v string) {
	w.Write([]byte{string_type})
	base64.NewEncoder(base64.StdEncoding, w).Write([]byte(v))
	w.Write([]byte{low})
}

func collate_slice(w io.Writer, v reflect.Value) {
	w.Write([]byte{slice_type})

	for i, l := 0, v.Len(); i < l; i++ {
		CollateValue(w, v.Index(i))
	}

	w.Write([]byte{low})
}

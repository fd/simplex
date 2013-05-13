package static

import (
	"reflect"
)

type G struct {
	t         Transformation
	tx        *Tx
	elem_type reflect.Type
	elems     map[string][]interface{}
}

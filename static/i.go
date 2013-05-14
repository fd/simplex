package static

import (
	"reflect"
	"simplex.sh/future"
)

type I struct {
	t         future.Deferred
	tx        *Tx
	elem_type reflect.Type
	elems     map[string]interface{}
}

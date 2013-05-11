package errors

import (
	"fmt"
)

type panic_error struct {
	v     interface{}
	stack []byte
}

func Panic(v interface{}, stack []byte) error {
	return &panic_error{v, stack}
}

func (e *panic_error) Error() string {
	return fmt.Sprintf("panic: %+v\n%s", e.v, e.stack)
}

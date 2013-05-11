package errors

import (
	"fmt"
)

func Fmt(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}

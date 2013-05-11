package errors

import (
	"fmt"
)

type forwarded_error struct {
	msg string
	err error
}

func Forward(err error, format string, args ...interface{}) error {
	return &forwarded_error{fmt.Sprintf(format, args...), err}
}

func (f *forwarded_error) Error() string {
	return fmt.Sprintf("%s (%s)", f.msg, f.err)
}

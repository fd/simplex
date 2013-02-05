package cas

import (
	"fmt"
)

type NotFound struct {
	Addr Addr
}

func (err NotFound) Error() string {
	return fmt.Sprintf("cas: Object not found (%s)", err.Addr)
}

func IsNotFound(err error) bool {
	_, ok := err.(NotFound)
	return err != nil && ok
}

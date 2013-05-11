package static

import (
	"fmt"
	"github.com/fd/static/errors"
	"reflect"
)

type Waiter interface {
	Wait() error
}

var waiterType = reflect.TypeOf((*Waiter)(nil)).Elem()

// Wait for all waiters in list
// list must be a slice of Waiters (or types witch implement Waiter)
func WaitForAll(list interface{}) error {
	if list == nil {
		return nil
	}

	var (
		err errors.List
	)

	v := reflect.ValueOf(list)
	t := v.Type()

	if t.Kind() != reflect.Slice || !t.Elem().Implements(waiterType) {
		panic(fmt.Sprintf("WaitForAll(list): expected []Waiter (%v %+v)", v.Type(), list))
	}

	for i, l := 0, v.Len(); i < l; i++ {
		w := v.Index(i).Interface().(Waiter)
		err.Add(w.Wait())
	}

	return err.Normalize()
}

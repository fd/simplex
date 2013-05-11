package static

import (
	"github.com/fd/static/errors"
	"runtime/debug"
	"sync"
)

type Transformation struct {
	wg  sync.WaitGroup
	err errors.List
}

func (t *Transformation) Wait() error {
	t.wg.Wait()
	return t.Err()
}

func (t *Transformation) Err() error {
	return t.err.Normalize()
}

func (t *Transformation) Do(f func() error) {
	t.wg.Add(1)
	go t.go_do(f)
}

func (t *Transformation) go_do(f func() error) {
	defer t.wg.Done()

	defer func() {
		r := recover()

		if r == nil {
			return
		}

		// if e, ok := r.(error); ok {
		//   t.err.Add(e)
		//   return
		// }

		t.err.Add(errors.Panic(r, debug.Stack()))
	}()

	err := f()
	if err != nil {
		t.err.Add(err)
	}
}

package future

import (
	"runtime/debug"
	"simplex.sh/errors"
	"sync"
)

type D interface {
	Wait() error
	Err() error
}

type Deferred struct {
	wg  sync.WaitGroup
	err errors.List
}

func (t *Deferred) Wait() error {
	t.wg.Wait()
	return t.Err()
}

func (t *Deferred) Err() error {
	return t.err.Normalize()
}

func (t *Deferred) Do(f func() error) {
	t.wg.Add(1)
	go t.go_do(f)
}

func (t *Deferred) go_do(f func() error) {
	defer t.wg.Done()

	defer func() {
		r := recover()

		if r == nil {
			return
		}

		t.err.Add(errors.Panic(r, debug.Stack()))
	}()

	err := f()

	if err != nil {
		t.err.Add(err)
	}
}

package runtime

import (
	"fmt"
	"simplex.sh/cas"
)

type (
	WorkerError struct {
		w      *worker_t
		data   interface{}
		err    error
		caller []byte
	}

	Changed struct {
		Id  [][]byte
		Key cas.Addr
		A   cas.Addr
		B   cas.Addr
	}
)

func (e *WorkerError) Event() string { return e.Error() }
func (e *WorkerError) Error() string { return fmt.Sprintf("%s: %s\n%s", e.w, e.err, e.caller) }

func (e *Changed) Event() string {
	return "Changed()"
}

func (e *Changed) Depth() int {
	return len(e.Id) - 1
}

func (e *Changed) IdAt(idx int) []byte {
	if idx < 0 {
		idx = len(e.Id) + idx
	}
	return e.Id[idx]
}

package runtime

import (
	"fmt"
	"github.com/fd/simplex/data/storage"
)

type (
	Event interface {
		isEvent()
	}

	EvError interface {
		Event
		Error() string
	}

	EvResolvedTable interface {
		Event
		Table() Table
	}

	ev_DONE_worker struct {
		w *worker_t
	}

	ev_DONE_pool struct {
		p *worker_pool_t
	}

	ev_ERROR struct {
		w      *worker_t
		data   interface{}
		err    error
		caller []byte
	}

	// a unit of progres from a -> b
	// representing a changing key/value
	// a is ZeroSHA when adding the key
	// b is ZeroSHA when remove the key
	ev_CHANGE struct {
		table string
		a     storage.SHA
		b     storage.SHA
	}

	// a unit of progres from a -> b
	// representing a changing table
	// a is ZeroSHA when adding the table
	// b is ZeroSHA when remove the table
	ev_CONSISTENT struct {
		table   string
		old_sha storage.SHA
		new_sha storage.SHA
	}
)

func (*ev_DONE_worker) isEvent() {}
func (*ev_DONE_pool) isEvent()   {}
func (*ev_ERROR) isEvent()       {}
func (*ev_CHANGE) isEvent()      {}
func (*ev_CONSISTENT) isEvent()  {}

func (e *ev_ERROR) Error() string { return fmt.Sprintf("%s\n%s", e.err, e.caller) }

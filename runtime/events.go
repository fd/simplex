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

	ev_SET struct {
		table   string
		old_sha storage.SHA
		new_sha storage.SHA
	}

	ev_DEL struct {
		table   string
		old_sha storage.SHA
	}

	ev_CONSISTENT struct {
		table   string
		old_sha storage.SHA
		new_sha storage.SHA
	}
)

func (*ev_DONE_worker) isEvent() {}
func (*ev_DONE_pool) isEvent()   {}
func (*ev_ERROR) isEvent()       {}
func (*ev_SET) isEvent()         {}
func (*ev_DEL) isEvent()         {}
func (*ev_CONSISTENT) isEvent()  {}

func (e *ev_ERROR) Error() string { return fmt.Sprintf("%s\n%s", e.err, e.caller) }

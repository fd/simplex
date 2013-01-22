package runtime

import (
	"fmt"
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
		old_sha string
		new_sha string
	}

	ev_DEL struct {
		table   string
		old_sha string
	}

	ev_CONSISTENT struct {
		table   string
		old_sha string
		new_sha string
	}
)

func (*ev_DONE_worker) isEvent() {}
func (*ev_DONE_pool) isEvent()   {}
func (*ev_ERROR) isEvent()       {}
func (*ev_SET) isEvent()         {}
func (*ev_DEL) isEvent()         {}
func (*ev_CONSISTENT) isEvent()  {}

func (e *ev_ERROR) Error() string { return fmt.Sprintf("%s\n%s", e.err, e.caller) }

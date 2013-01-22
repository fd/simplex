package runtime

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
		w    *worker_t
		data interface{}
	}
)

func (*ev_DONE_worker) isEvent() {}
func (*ev_DONE_pool) isEvent()   {}
func (*ev_ERROR) isEvent()       {}

func (*ev_ERROR) Error() string { return "(error)" }

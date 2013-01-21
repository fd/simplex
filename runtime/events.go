package runtime

type Event interface {
	isEvent()
}

type ev_done struct {
	w *worker_t
}

type ev_error struct {
	w    *worker_t
	data interface{}
}

func (*ev_done) isEvent()  {}
func (*ev_error) isEvent() {}

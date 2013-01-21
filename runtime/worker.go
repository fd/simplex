package runtime

type worker_op_k uint

type worker_op_t struct {
	kind worker_op_k
	data interface{}
}

type worker_t struct {
	txn         *Transaction
	def         Deferred
	subscribers []chan<- Event
	operations  chan worker_op_t
}

const (
	op_SUB worker_op_k = iota
)

func (w *worker_t) run(worker_events chan<- Event) {
	events := make(chan Event, 2)
	w.operations = make(chan worker_op_t, 2)

	go w.go_resolve(events)
	go w.go_run(events, worker_events)
}

func (w *worker_t) go_resolve(events chan<- Event) {
	defer func() {
		if e := recover(); e != nil {
			events <- &ev_error{w, e}
		}
		events <- &ev_done{w}
		close(events)
	}()

	w.def.Resolve(w.txn, events)
}

func (w *worker_t) go_run(events <-chan Event, worker_events chan<- Event) {
	log := make([]Event, 0, 128)

	defer func() {
		if e := recover(); e != nil {
			worker_events <- &ev_error{w, e}
		}
		worker_events <- &ev_done{w}
	}()

	for {
		select {

		case e := <-events:

			log = append(log, e)

			for _, sub := range w.subscribers {
				sub <- e
			}

			if _, ok := e.(*ev_error); ok {
				worker_events <- e
			}

			if _, ok := e.(*ev_done); ok {
				return
			}

		case op := <-w.operations:
			switch op.kind {
			case op_SUB:
				if ch, ok := op.data.(chan<- Event); ok {
					for _, e := range log {
						ch <- e
					}
					w.subscribers = append(w.subscribers, ch)
				}

			default:
				panic("not reached")
			}

		}
	}
}

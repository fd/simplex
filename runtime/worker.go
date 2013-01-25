package runtime

import (
	"fmt"
	"runtime/debug"
)

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

func (w *worker_t) subscribe() <-chan Event {
	ch := make(chan Event, 1)
	w.operations <- worker_op_t{op_SUB, ch}
	return ch
}

func (w *worker_t) run(worker_events chan<- Event) {
	events := make(chan Event, 2)
	w.operations = make(chan worker_op_t, 2)

	go w.go_resolve(events)
	go w.go_run(events, worker_events)
}

func (w *worker_t) go_resolve(events chan<- Event) {
	defer func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				events <- &ev_ERROR{w, e, err, debug.Stack()}
			} else {
				events <- &ev_ERROR{w, e, fmt.Errorf("error: %+v", e), debug.Stack()}
			}
		}
		events <- &ev_DONE_worker{w}
		close(events)
	}()

	w.def.Resolve(w.txn, events)
}

func (w *worker_t) go_run(events <-chan Event, worker_events chan<- Event) {
	log := make([]Event, 0, 128)

	defer func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				worker_events <- &ev_ERROR{w, e, err, debug.Stack()}
			} else {
				worker_events <- &ev_ERROR{w, e, fmt.Errorf("error: %+v", e), debug.Stack()}
			}
		}
		worker_events <- &ev_DONE_worker{w}

		for _, sub := range w.subscribers {
			close(sub)
		}
	}()

	for {
		select {

		case e := <-events:

			log = append(log, e)

			for _, sub := range w.subscribers {
				sub <- e
			}

			if _, ok := e.(*ev_CONSISTENT); ok {
				worker_events <- e
			}

			if _, ok := e.(*ev_ERROR); ok {
				worker_events <- e
			}

			if _, ok := e.(*ev_DONE_worker); ok {
				return
			}

		case op := <-w.operations:
			switch op.kind {
			case op_SUB:
				if ch, ok := op.data.(chan Event); ok {
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

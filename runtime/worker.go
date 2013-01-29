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
	op_DONE
)

func (w *worker_t) String() string {
	return "Worker(" + w.def.DeferredId() + ")"
}

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
	var (
		log            = make([]Event, 0, 128)
		worker_is_done bool
	)

	defer func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				worker_events <- &ev_ERROR{w, e, err, debug.Stack()}
			} else {
				worker_events <- &ev_ERROR{w, e, fmt.Errorf("error: %+v", e), debug.Stack()}
			}
		}
		if !worker_is_done {
			for _, sub := range w.subscribers {
				close(sub)
			}
			worker_events <- &ev_DONE_worker{w}
		}
	}()

	for {
		select {

		case e := <-events:
			log, worker_is_done = w.handle_event(e, log, worker_events)

		case op := <-w.operations:
			exit := w.handle_operation(op, log, worker_is_done)
			if exit {
				return
			}

		}

		if worker_is_done {
			break
		}
	}

	for op := range w.operations {
		exit := w.handle_operation(op, log, worker_is_done)
		if exit {
			return
		}
	}
}

func (w *worker_t) handle_event(e Event, log []Event, worker_events chan<- Event) (log_o []Event, done bool) {
	log_o = log
	done = false

	switch e.(type) {

	case nil:
		// ignore

	case *EvConsistent:
		log_o = append(log_o, e)
		for _, sub := range w.subscribers {
			sub <- e
		}
		worker_events <- e

	case *ev_ERROR:
		worker_events <- e

	case *ev_DONE_worker:
		done = true
		for _, sub := range w.subscribers {
			close(sub)
		}
		worker_events <- &ev_DONE_worker{w}

	default:
		log_o = append(log_o, e)
		for _, sub := range w.subscribers {
			sub <- e
		}
		//worker_events <- e

	}

	return
}

func (w *worker_t) handle_operation(op worker_op_t, log []Event, done bool) (exit bool) {
	switch op.kind {

	case op_SUB:
		if ch, ok := op.data.(chan Event); ok {
			for _, e := range log {
				ch <- e
			}
			if done {
				close(ch)
			} else {
				w.subscribers = append(w.subscribers, ch)
			}
		}

	case op_DONE:
		exit = true

	default:
		panic("not reached")

	}

	return
}

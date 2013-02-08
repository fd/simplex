package runtime

import (
	"fmt"
	"github.com/fd/simplex/runtime/event"
	"runtime/debug"
	"sync"
)

type worker_t struct {
	txn *Transaction
	def Deferred
}

func (w *worker_t) String() string {
	return "Worker(" + w.def.DeferredId() + ")"
}

func (w *worker_t) run(wg *sync.WaitGroup) {
	events := w.txn.dispatcher.Register(w.def.DeferredId())

	go w.go_resolve(events, wg)
}

func (w *worker_t) go_resolve(events chan<- event.Event, wg *sync.WaitGroup) {
	defer func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				events <- &WorkerError{w, e, err, debug.Stack()}
			} else {
				events <- &WorkerError{w, e, fmt.Errorf("error: %+v", e), debug.Stack()}
			}
		}
		close(events)
		wg.Done()
	}()

	w.def.Resolve(w.txn, events)
}

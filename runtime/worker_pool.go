package runtime

import (
	"fmt"
)

type worker_pool_t struct {
	scheduled_workers chan *schedule_worker_op
}

type schedule_worker_op struct {
	def   Deferred
	txn   *Transaction
	reply chan *worker_t
}

func (p *worker_pool_t) run() <-chan Event {
	done := make(chan Event, 1)
	p.scheduled_workers = make(chan *schedule_worker_op, 1)

	go p.go_run(done)

	return done
}

func (p *worker_pool_t) go_run(events chan<- Event) {
	var (
		worker_events = make(chan Event, 1)
		workers       = map[*worker_t]bool{}
		deferreds     = map[string]*worker_t{}
	)

	defer func() {
		for _, w := range deferreds {
			w.operations <- worker_op_t{op_DONE, nil}
		}
		events <- &ev_DONE_pool{p}
		close(events)
	}()

	for {
		select {

		case e := <-worker_events:
			events <- e

			if done, ok := e.(*ev_DONE_worker); ok {
				delete(workers, done.w)
				if len(workers) == 0 {
					return
				}
			}

		case op := <-p.scheduled_workers:
			w, ok := deferreds[op.def.DeferredId()]
			if !ok {
				w = &worker_t{def: op.def, txn: op.txn}
				workers[w] = true
				deferreds[op.def.DeferredId()] = w
				w.run(worker_events)
				fmt.Println("ADD:", w)
			} else {
				fmt.Println("SUB:", w)
			}
			op.reply <- w
			close(op.reply)

		}
	}
}

func (p *worker_pool_t) schedule(txn *Transaction, def Deferred) *worker_t {
	op := &schedule_worker_op{
		txn:   txn,
		def:   def,
		reply: make(chan *worker_t, 1),
	}

	p.scheduled_workers <- op
	return <-op.reply
}

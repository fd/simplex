package runtime

type worker_pool_t struct {
	workers           map[*worker_t]bool
	scheduled_workers chan *worker_t
}

func (p *worker_pool_t) run() <-chan Event {
	done := make(chan Event, 1)
	p.scheduled_workers = make(chan *worker_t, 1)
	p.workers = make(map[*worker_t]bool)

	go p.go_run(done)

	return done
}

func (p *worker_pool_t) go_run(events chan<- Event) {
	defer func() {
		events <- &ev_DONE_pool{p}
		close(events)
	}()

	worker_events := make(chan Event, 1)

	for {
		select {

		case e := <-worker_events:
			if done, ok := e.(*ev_DONE_worker); ok {
				delete(p.workers, done.w)
				if len(p.workers) == 0 {
					return
				}
			}

			events <- e

		case w := <-p.scheduled_workers:
			p.workers[w] = true
			w.run(worker_events)

		}
	}
}

func (p *worker_pool_t) schedule(txn *Transaction, def Deferred) {
	p.scheduled_workers <- &worker_t{txn: txn, def: def}
}

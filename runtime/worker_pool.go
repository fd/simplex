package runtime

import (
	"fmt"
	"sync"
)

type worker_pool_t struct {
	operations chan *schedule_worker_op
}

type schedule_worker_op struct {
	def Deferred
	txn *Transaction
}

func (p *worker_pool_t) run() <-chan bool {
	done := make(chan bool, 1)
	p.operations = make(chan *schedule_worker_op, 1)

	go p.go_run(done)

	return done
}

func (p *worker_pool_t) go_run(done chan<- bool) {
	var (
		workers = map[string]*worker_t{}
		wg      sync.WaitGroup
	)

	defer func() {
		done <- true
		close(done)
	}()

	go func() {
		wg.Wait()
		close(p.operations)
	}()

	for op := range p.operations {

		w := workers[op.def.DeferredId()]
		if w == nil {
			w = &worker_t{def: op.def, txn: op.txn}
			workers[op.def.DeferredId()] = w

			wg.Add(1)
			w.run(&wg)

			fmt.Println("ADD:", w)
		} else {
			fmt.Println("SUB:", w)
		}

	}
}

func (p *worker_pool_t) schedule(txn *Transaction, def Deferred) {
	op := &schedule_worker_op{
		txn: txn,
		def: def,
	}

	p.operations <- op
}

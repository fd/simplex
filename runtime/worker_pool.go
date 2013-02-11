package runtime

import (
	"sync"
)

type worker_pool_t struct {
	operations chan *schedule_worker_op
	wg         *sync.WaitGroup
}

type schedule_worker_op struct {
	def   Deferred
	txn   *Transaction
	reply chan bool
}

func (p *worker_pool_t) Start() {
	p.operations = make(chan *schedule_worker_op, 1)
	p.wg = &sync.WaitGroup{}

	go p.go_run()
}

func (p *worker_pool_t) Stop() {
	p.wg.Wait()
	close(p.operations)
}

func (p *worker_pool_t) go_run() {
	var (
		workers = map[string]bool{}
		wg      = p.wg
	)

	for op := range p.operations {

		started := workers[op.def.DeferredId()]
		if !started {
			w := &worker_t{def: op.def, txn: op.txn}
			workers[op.def.DeferredId()] = true

			wg.Add(1)
			w.run(wg)
		}
		op.reply <- true
		close(op.reply)

	}
}

func (p *worker_pool_t) schedule(txn *Transaction, def Deferred) {
	reply := make(chan bool, 1)
	p.operations <- &schedule_worker_op{
		txn:   txn,
		def:   def,
		reply: reply,
	}
	<-reply
}

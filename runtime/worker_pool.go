package runtime

import (
	"sync"
)

type worker_pool_t struct {
	operations chan *schedule_worker_op
	wg         sync.WaitGroup
}

type schedule_worker_op struct {
	def   Deferred
	txn   *Transaction
	reply chan bool
}

func (p *worker_pool_t) Start() {
	p.operations = make(chan *schedule_worker_op, 1)

	go p.go_run()
}

func (p *worker_pool_t) Stop() {
	p.wg.Wait()
	close(p.operations)
}

func (p *worker_pool_t) go_run() {
	var (
		workers = map[string]*worker_t{}
	)

	for op := range p.operations {

		w := workers[op.def.DeferredId()]
		if w == nil {
			w = &worker_t{def: op.def, txn: op.txn}
			workers[op.def.DeferredId()] = w

			p.wg.Add(1)
			w.run(&p.wg)
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

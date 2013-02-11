package event

import (
	"runtime"
	"sync"
)

type (
	Dispatcher struct {
		operations chan interface{}
		exchanges  map[string]*exchange
	}

	Subscription struct {
		C <-chan Event

		exchange *exchange
		outbound chan Event
		cursor   int
	}

	exchange struct {
		name    string
		broken  bool
		inbound chan Event
		log     []Event
		rw_mtx  *sync.RWMutex
		cond    *sync.Cond
	}

	disp_op__register struct {
		name  string
		reply chan chan<- Event
	}

	disp_op__subscribe struct {
		name  string
		reply chan *exchange
	}

	disp_op__stop struct {
	}
)

func (disp *Dispatcher) Start() {
	if disp.operations == nil {
		disp.operations = make(chan interface{}, 1)
		disp.exchanges = make(map[string]*exchange)
		go disp.go_run()
	}
}

func (disp *Dispatcher) Stop() {
	disp.operations <- &disp_op__stop{}
}

func (disp *Dispatcher) Subscribe(name string) *Subscription {
	reply := make(chan *exchange, 1)
	disp.operations <- &disp_op__subscribe{name, reply}
	exch := <-reply

	out := make(chan Event, 1)
	sub := &Subscription{
		C:        out,
		outbound: out,
		exchange: exch,
	}

	go sub.go_run()

	return sub
}

// Returns a named channel.
func (disp *Dispatcher) Register(name string) chan<- Event {
	reply := make(chan chan<- Event, 1)
	disp.operations <- &disp_op__register{name, reply}
	return <-reply
}

func (disp *Dispatcher) register(name string) *exchange {
	e := disp.exchanges[name]
	if e != nil {
		return e
	}

	rw_mtx := &sync.RWMutex{}
	e = &exchange{
		name:    name,
		inbound: make(chan Event, 1),
		rw_mtx:  rw_mtx,
		cond:    sync.NewCond(rw_mtx),
	}

	disp.exchanges[name] = e

	go e.go_run()

	return e
}

func (disp *Dispatcher) go_run() {
	for op := range disp.operations {
		switch o := op.(type) {

		case *disp_op__register:
			o.reply <- disp.register(o.name).inbound
			close(o.reply)

		case *disp_op__subscribe:
			e := disp.register(o.name)
			o.reply <- e

		case *disp_op__stop:
			close(disp.operations)
			for _, e := range disp.exchanges {
				e.break_exchange()
			}
			return

		}
	}
}

func (sub *Subscription) go_run() {
	for {
		closed, broken := sub.pop_events()

		if closed || broken {
			break
		}
	}

	ensure_closed(sub.outbound)
}

func (sub *Subscription) get_log() (log []Event, closed, broken bool) {
	sub.exchange.cond.L.Lock()
	defer sub.exchange.cond.L.Unlock()
	// wait for event
	for sub.cursor == len(sub.exchange.log) && sub.exchange.inbound != nil && !sub.exchange.broken {
		sub.exchange.cond.Wait()
	}

	return sub.exchange.log, sub.exchange.inbound == nil, sub.exchange.broken
}

func (sub *Subscription) pop_events() (closed, broken bool) {
	var (
		log []Event
	)

	log, closed, broken = sub.get_log()

	for sub.cursor < len(log) {
		runtime.Gosched()
		event := log[sub.cursor]
		if closed || broken {
			sub.outbound <- event
			sub.cursor += 1
		} else {
			select {
			case sub.outbound <- event:
				sub.cursor += 1
			default:
				closed = false
				broken = false
				return
			}
		}
	}

	return
}

func (exch *exchange) push_event(event Event) {
	exch.cond.L.Lock()
	defer exch.cond.L.Unlock()

	exch.log = append(exch.log, event)

	exch.cond.Broadcast()
}

func (exch *exchange) break_exchange() {
	exch.cond.L.Lock()
	defer exch.cond.L.Unlock()

	ensure_closed(exch.inbound)
	exch.inbound = nil // closed
	exch.broken = true

	exch.cond.Broadcast()
}

func (exch *exchange) close_exchange() {
	exch.cond.L.Lock()
	defer exch.cond.L.Unlock()

	ensure_closed(exch.inbound)
	exch.inbound = nil // closed

	exch.cond.Broadcast()
}

func (exch *exchange) go_run() {
	for event := range exch.inbound {
		exch.push_event(event)
	}

	exch.close_exchange()
}

func ensure_closed(c chan Event) {
	defer func() { recover() }()
	close(c)
}

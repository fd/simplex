package event

import (
	"fmt"
)

type (
	Dispatcher struct {
		operations chan interface{}
		exchanges  map[string]*exchange
		curr_id    int
	}

	Subscription struct {
		C    <-chan Event
		name string
		id   int
		c    chan Event
		disp *Dispatcher
	}

	exchange struct {
		name string
		c    chan Event
	}

	disp_op__register struct {
		name  string
		reply chan chan<- Event
	}

	disp_op__subscribe struct {
		name  string
		reply chan Subscription
	}

	disp_op__cancel struct {
		name  string
		id    int
		reply chan bool
	}

	disp_op__stop struct {
	}

	exch_op__subscribe struct {
		id int
		c  chan Event
	}

	exch_op__cancel struct {
		id int
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

func (disp *Dispatcher) Subscribe(name string) Subscription {
	reply := make(chan Subscription, 1)
	disp.operations <- &disp_op__subscribe{name, reply}
	return <-reply
}

// Returns a named channel.
func (disp *Dispatcher) Register(name string) chan<- Event {
	reply := make(chan chan<- Event, 1)
	disp.operations <- &disp_op__register{name, reply}
	return <-reply
}

func (sub Subscription) Cancel() {
	defer func() { recover() }()
	ensure_closed(sub.c)
	reply := make(chan bool, 1)
	sub.disp.operations <- &disp_op__cancel{sub.name, sub.id, reply}
	<-reply
}

func (disp *Dispatcher) register(name string) *exchange {
	e := disp.exchanges[name]
	if e != nil {
		return e
	}

	e = &exchange{
		name: name,
		c:    make(chan Event, 1),
	}

	disp.exchanges[name] = e

	go e.go_run()

	return e
}

func (disp *Dispatcher) go_run() {
	for op := range disp.operations {
		switch o := op.(type) {

		case *disp_op__register:
			fmt.Printf("pubsub: REGISTER  %s\n", o.name)
			o.reply <- disp.register(o.name).c
			close(o.reply)

		case *disp_op__subscribe:
			fmt.Printf("pubsub: SUBSCRIBE %s\n", o.name)
			disp.curr_id += 1
			c := make(chan Event, 1)

			disp.register(o.name).send_op(&exch_op__subscribe{disp.curr_id, c})
			o.reply <- Subscription{c, o.name, disp.curr_id, c, disp}

		case *disp_op__cancel:
			disp.register(o.name).send_op(&exch_op__cancel{o.id})
			o.reply <- true

		case *disp_op__stop:
			close(disp.operations)
			for _, e := range disp.exchanges {
				ensure_closed(e.c)
			}
			return

		}
	}
}

func ensure_closed(c chan Event) {
	defer func() { recover() }()
	close(c)
}

func try_send(c chan<- Event, e Event) {
	defer func() { recover() }()
	c <- e
}

func (exch *exchange) send_op(op Event) {
	defer func() {
		// no error
		if e := recover(); e == nil {
			return
		} else {
			fmt.Printf("ERR: %s\n", e)
		}

		// channel was closed, make a new one
		exch.c = make(chan Event, 1)
		exch.c <- op
	}()

	exch.c <- op
}

func (exch *exchange) go_run() {
	var (
		subscribers map[int]chan Event
		log         []Event
	)

	subscribers = make(map[int]chan Event, 1)

	for event := range exch.c {
		switch e := event.(type) {

		case *exch_op__subscribe:
			// add to subscribers
			// deliver event log
			subscribers[e.id] = e.c
			for _, log_e := range log {
				e.c <- log_e
			}

		case *exch_op__cancel:
			// remove from subscribers
			// close channel
			if c, ok := subscribers[e.id]; ok {
				ensure_closed(c)
				delete(subscribers, e.id)
			}

		default:
			// log event
			// deliver event to subscriber
			log = append(log, e)
			for _, sub := range subscribers {
				try_send(sub, e)
			}

		}
	}

	// close existing subscribers
	for _, sub := range subscribers {
		close(sub)
	}
	subscribers = make(map[int]chan Event, 1)

	fmt.Printf("pubsub: LOGGING %s\n", exch.name)

	// go in log delivery mode
	for event := range exch.c {
		switch op := event.(type) {

		case *exch_op__subscribe:
			// add to subscribers
			// deliver event log
			subscribers[op.id] = op.c
			go go_deliver_log(op.c, log)

		case *exch_op__cancel:
			// remove from subscribers
			// close channel
			if c, ok := subscribers[op.id]; ok {
				ensure_closed(c)
				delete(subscribers, op.id)
			}
		}
	}

	// close existing subscribers
	for _, sub := range subscribers {
		close(sub)
	}
	subscribers = nil

	fmt.Printf("pubsub: CLOSED %s\n", exch.name)
}

func go_deliver_log(c chan<- Event, log []Event) {
	defer close(c)
	defer func() {
		e := recover()
		if e != nil {
			fmt.Printf("ERR: %s\n", e)
		}
	}()
	for _, log_e := range log {
		c <- log_e
	}
}

func (*exch_op__cancel) Event() string    { return "[INTERNAL: exch_op__cancel]" }
func (*exch_op__subscribe) Event() string { return "[INTERNAL: exch_op__subscribe]" }

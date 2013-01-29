package runtime

import (
	"sync"
)

type (
	Funnel struct {
		inbound   []<-chan Event
		outbound  <-chan Event
		collector chan Event
	}
)

func (f *Funnel) Add(ch <-chan Event) {
	f.inbound = append(f.inbound, ch)
}

func (f *Funnel) Run() <-chan Event {
	if f.outbound != nil {
		return f.outbound
	}

	if len(f.inbound) == 0 {
		f.collector = make(chan Event, 1)
		f.outbound = f.collector
		close(f.collector)
		return f.outbound
	}

	if len(f.inbound) == 1 {
		f.outbound = f.inbound[0]
		return f.outbound
	}

	f.collector = make(chan Event, 1)
	f.outbound = f.collector

	go f.go_sink()

	return f.outbound
}

func (f *Funnel) go_sink() {
	var wg sync.WaitGroup
	wg.Add(len(f.inbound))

	defer close(f.collector)

	for _, ch := range f.inbound {
		go f.go_collect(&wg, ch)
	}

	wg.Wait()
}

func (f *Funnel) go_collect(wg *sync.WaitGroup, ch <-chan Event) {
	defer wg.Done()

	for e := range ch {
		f.collector <- e
	}
}
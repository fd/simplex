package event

import (
	"sync"
)

type (
	Funnel struct {
		inbound  []<-chan Event
		outbound <-chan Event
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
		collector := make(chan Event, 1)
		f.outbound = collector
		close(collector)
		return f.outbound
	}

	if len(f.inbound) == 1 {
		f.outbound = f.inbound[0]
		return f.outbound
	}

	collector := make(chan Event, 1)
	f.outbound = collector

	go f.go_sink(collector)

	return f.outbound
}

func (f *Funnel) go_sink(collector chan Event) {
	var wg sync.WaitGroup
	wg.Add(len(f.inbound))

	defer close(collector)

	for _, ch := range f.inbound {
		go f.go_collect(&wg, collector, ch)
	}

	wg.Wait()
}

func (f *Funnel) go_collect(wg *sync.WaitGroup, collector chan Event, ch <-chan Event) {
	defer wg.Done()

	for e := range ch {
		collector <- e
	}
}

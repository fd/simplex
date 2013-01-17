package data

import (
	"fmt"
	"github.com/fd/simplex/w/data/storage"
	"github.com/fd/simplex/w/util"
)

var current_engine = NewEngine()

func Setup(store_url string) error {
	return current_engine.Setup(store_url)
}

func Update(c Changes) {
	current_engine.Update(c)
}

func Reset() {
	current_engine.Reset()
}

func Stop() {
	current_engine.Stop()
}

func Wait() {
	current_engine.Wait()
}

func Run() {
	current_engine.Run()
}

type Engine struct {
	store *storage.S

	transactions chan *transaction
	done         chan bool

	transformations         map[string]transformation
	transformation_counters map[string]int
}

func NewEngine() *Engine {
	return &Engine{
		transformations:         make(map[string]transformation),
		transformation_counters: make(map[string]int),
		transactions:            make(chan *transaction),
		done:                    make(chan bool, 1),
	}
}

func (e *Engine) Setup(store_url string) error {
	s, err := storage.New(store_url)
	if err != nil {
		return err
	}

	e.store = s
	return nil
}

func (e *Engine) Update(changes Changes) {
	txn := new_transaction(e, changes)
	e.schedule(txn)
}

func (e *Engine) Reset() {
}

func (e *Engine) Stop() {
	close(e.transactions)
	e.Wait()
}

func (e *Engine) Wait() {
	<-e.done
	e.done <- true
}

func (e *Engine) Run() {
	go e.go_run()
}

func (e *Engine) go_run() {
	for txn := range e.transactions {
		txn.project()
	}

	e.done <- true
}

func (e *Engine) UnscopedView() View {
	return View{engine: e, current: nil}
}

func (e *Engine) ScopedView() View {
	v := e.UnscopedView()
	app := util.InitializingApplication()

	fmt.Printf("ScopedView[%s]\n", app)

	if app == "unknown" {
		panic("Initialized view in unknown application")
	}

	return v
}

func (e *Engine) sorted_transformations() []transformation {
	present := map[string]bool{}
	transformations := make([]transformation, 0, len(e.transformations))

	for id, transformation := range e.transformations {
		if present[id] {
			continue
		}

		for _, dep := range transformation.Dependencies() {
			if present[dep.Id()] {
				continue
			}

			transformations = append(transformations, dep)
			present[dep.Id()] = true
		}

		transformations = append(transformations, transformation)
		present[id] = true
	}

	return transformations
}

func (e *Engine) schedule(txn *transaction) {
	e.transactions <- txn
}

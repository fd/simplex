package data

import (
	"fmt"
	"github.com/fd/w/util"
)

var current_engine = NewEngine()

type Engine struct {
	/*  source *Source
	target *Target
	state  *State*/

	transformations         map[string]transformation
	transformation_counters map[string]int
}

func NewEngine() *Engine {
	return &Engine{
		transformations:         make(map[string]transformation),
		transformation_counters: make(map[string]int),
	}
}

func (e *Engine) Update(changes Changes) {
	txn := new_transaction(e, changes)
	e.schedule(txn)
}

func (e *Engine) Reset() {
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
}

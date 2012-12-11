package data

import (
	"fmt"
	"github.com/fd/w/util"
)

var current_engine = &Engine{}

type Engine struct {
	source *Source
	target *Target
	state  *State
}

type Changes struct {
	Create  []Value
	Update  map[string]Value
	Destroy []string
}

func (e *Engine) Update(changes Changes) {
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

type StoreReader interface {
	Ids() []string
	Get(id string) Value
}

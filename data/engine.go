package data

type Engine struct {
	source Source
	target Target
	state  State

	transformations map[string]*transformation_controller
}

func (e *Engine) Update(changes Changes) {
}

func (e *Engine) Reset() {
	e.state.Flush()
	e.Update(&Changes{Added: e.source.Ids()})
}

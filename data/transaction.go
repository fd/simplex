package data

type transaction struct {
	engine          *Engine
	changes         Changes
	transformations []transformation

	upstream_states map[string][]upstream_state
}

func new_transaction(e *Engine, changes Changes) *transaction {
	transformations := e.sorted_transformations()
	changes.engine = e

	return &transaction{
		engine:          e,
		changes:         changes,
		transformations: transformations,
		upstream_states: make(map[string][]upstream_state, len(transformations)),
	}
}

func (txn *transaction) Restore(s *state) {
}

func (txn *transaction) Save(s *state) {
}

func (txn *transaction) Propagate(ts []transformation, s *state) {
	for _, transformation := range ts {
		a := txn.upstream_states[transformation.Id()]
		a = append(a, s)
		txn.upstream_states[transformation.Id()] = a
	}
}

func (txn *transaction) project() {
	// bind tips
	for _, t := range txn.transformations {
		if len(t.Chain()) == 1 && len(txn.upstream_states[t.Id()]) == 0 {
			txn.upstream_states[t.Id()] = []upstream_state{txn.changes}
		}
	}

	// transform
	for _, transformation := range txn.transformations {
		txn.transform(transformation)
	}
}

func (txn *transaction) transform(t transformation) {
	states := txn.upstream_states[t.Id()]

	for _, state := range states {
		t.Transform(state, txn)
	}
}

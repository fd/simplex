package data

import (
	"fmt"
	"strings"
	"time"
)

type transaction struct {
	engine          *Engine
	changes         Changes
	transformations []transformation

	upstream_states map[string][]upstream_state

	saved_states map[string]interface{}
}

func new_transaction(e *Engine, changes Changes) *transaction {
	transformations := e.sorted_transformations()
	changes.engine = e

	return &transaction{
		engine:          e,
		changes:         changes,
		transformations: transformations,
		upstream_states: make(map[string][]upstream_state, len(transformations)),
		saved_states:    make(map[string]interface{}),
	}
}

func (txn *transaction) Restore(s *state, info interface{}) {
	txn.engine.state_table.Restore(strings.Join(s.Id(), "/"), info)
}

func (txn *transaction) Save(s *state) {
	txn.saved_states[strings.Join(s.Id(), "/")] = s.Info
}

func (txn *transaction) Commit() {
	{
		set := map[string]Value{}
		del := []string{}

		for id, val := range txn.changes.Create {
			set[id] = val
		}

		for id, val := range txn.changes.Update {
			set[id] = val
		}

		for _, id := range txn.changes.Destroy {
			del = append(del, id)
		}

		txn.engine.source_table.Commit(set, del)
	}

	{
		set := map[string]Value{}
		del := []string{}

		for id, val := range txn.saved_states {
			set[id] = val
		}

		txn.engine.state_table.Commit(set, del)
	}
}

func (txn *transaction) Propagate(ts []transformation, s *state) {
	for _, transformation := range ts {
		a := txn.upstream_states[transformation.Id()]
		a = append(a, s)
		txn.upstream_states[transformation.Id()] = a
	}
}

func (txn *transaction) project() {
	now := time.Now()

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

	fmt.Printf("txn: %f\n", float64(time.Now().Sub(now).Nanoseconds())/1000000)
	txn.Commit()
	fmt.Printf("txn: %f\n", float64(time.Now().Sub(now).Nanoseconds())/1000000)
}

func (txn *transaction) transform(t transformation) {
	states := txn.upstream_states[t.Id()]

	for _, state := range states {
		now := time.Now()

		c := len(state.Added()) + len(state.Changed()) + len(state.Removed())
		fmt.Printf("beg A=>B: %s\n    %s\n    %d\n", t.Id(), state.Id(), c)

		t.Transform(state, txn)

		fmt.Printf("end A=>B: %f\n", float64(time.Now().Sub(now).Nanoseconds())/1000000)
	}
}

package runtime

import (
	"github.com/fd/simplex/runtime/event"
)

/*
  Void() registers a side-effect free terminal. It is mainly useful for debugging
  as it ensurs that the Deferred def is resolved.
*/
func Void(def Deferred) {
	Env.RegisterTerminal(&void_terminal{def})
}

type void_terminal struct {
	def Deferred
}

func (t *void_terminal) DeferredId() string {
	return "void(" + t.def.DeferredId() + ")"
}

func (t *void_terminal) Resolve(txn *Transaction, events chan<- event.Event) {
	for e := range txn.Resolve(t.def).C {
		// propagate error events
		if err, ok := e.(event.Error); ok {
			events <- err
			continue
		}

		// ignore
	}
}

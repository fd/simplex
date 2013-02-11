package runtime

import (
	"simplex.sh/runtime/event"
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
	src_events := txn.Resolve(t.def)

	for e := range src_events.C {
		// propagate error events
		if err, ok := e.(event.Error); ok {
			events <- err
			continue
		}

		// ignore
	}
}

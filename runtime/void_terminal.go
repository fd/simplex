package runtime

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

func (t *void_terminal) Resolve(txn *Transaction, events chan<- Event) {
	t.def.Resolve(txn, events)
}

package runtime

/*
  Void() registers a side-effect free terminal. It is mainly useful for debugging
  as it ensurs that the Deferred def is resolved.
*/
func Void(r Resolver) {
	Env.RegisterTerminal(&void_terminal{r})
}

type void_terminal struct {
	r Resolver
}

func (t *void_terminal) DeferredId() string {
	return "void(" + t.r.DeferredId() + ")"
}

func (t *void_terminal) Resolve(state *Transaction) IChange {
	state.Resolve(t.r)
	return IChange{}
}

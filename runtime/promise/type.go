package promise

import (
	"simplex.sh/cas"
	"simplex.sh/cas/btree"
	"simplex.sh/runtime/event"
)

type (
	Deferred interface {
		DeferredId() string
		Resolve(state State, events chan<- event.Event)
	}

	State interface {
		Store() cas.Store

		Resolve(Deferred) *event.Subscription

		GetTable(string) *btree.Tree
		CommitTable(string, *btree.Tree) (prev, curr cas.Addr)
	}
)

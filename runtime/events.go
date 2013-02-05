package runtime

import (
	"fmt"
	"github.com/fd/simplex/cas"
	"github.com/fd/simplex/cas/btree"
)

type (
	Event interface {
		isEvent()
	}

	EvError interface {
		Event
		Error() string
	}

	ev_DONE_worker struct {
		w *worker_t
	}

	ev_DONE_pool struct {
		p *worker_pool_t
	}

	ev_ERROR struct {
		w      *worker_t
		data   interface{}
		err    error
		caller []byte
	}

	// a unit of progres from a -> b
	// representing a changing key/value
	// a is ZeroSHA when adding the key
	// b is ZeroSHA when remove the key
	ev_CHANGE struct {
		table        string
		collated_key []byte
		key          cas.Addr
		a            cas.Addr
		b            cas.Addr
	}

	// a unit of progres from a -> b
	// representing a changing table
	// a is ZeroSHA when adding the table
	// b is ZeroSHA when remove the table
	EvConsistent struct {
		Table string
		A     cas.Addr
		B     cas.Addr
	}
)

func (*ev_DONE_worker) isEvent() {}
func (*ev_DONE_pool) isEvent()   {}
func (*ev_ERROR) isEvent()       {}
func (*ev_CHANGE) isEvent()      {}
func (*EvConsistent) isEvent()   {}

func (e *ev_ERROR) Error() string { return fmt.Sprintf("%s: %s\n%s", e.w, e.err, e.caller) }

func (e *ev_DONE_worker) String() string {
	return "DONE(" + e.w.String() + ")"
}

func (e *EvConsistent) GetTableA(txn *Transaction) *btree.Tree {
	return txn.env.LoadTable(e.A)
}

func (e *EvConsistent) GetTableB(txn *Transaction) *btree.Tree {
	return txn.env.LoadTable(e.B)
}

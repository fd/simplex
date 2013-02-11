package runtime

import (
	"fmt"
	"simplex.sh/cas"
	"simplex.sh/cas/btree"
)

type (
	ev_DONE_pool struct {
		p *worker_pool_t
	}

	WorkerError struct {
		w      *worker_t
		data   interface{}
		err    error
		caller []byte
	}

	// a unit of progres from a -> b
	// representing a changing key/value
	// a is ZeroSHA when adding the key
	// b is ZeroSHA when remove the key
	ChangedMember struct {
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
	ConsistentTable struct {
		Table string
		A     cas.Addr
		B     cas.Addr
	}
)

func (*ev_DONE_pool) isEvent()       {}
func (e *WorkerError) Event() string { return e.Error() }
func (e *ChangedMember) Event() string {
	return fmt.Sprintf(
		"ChangedMember(table: %s, member: %s)",
		e.table, e.collated_key,
	)
}
func (e *ConsistentTable) Event() string {
	return fmt.Sprintf(
		"Consistent(table: %s)",
		e.Table,
	)
}

func (e *WorkerError) Error() string { return fmt.Sprintf("%s: %s\n%s", e.w, e.err, e.caller) }

func (e *ConsistentTable) GetTableA(txn *Transaction) *btree.Tree {
	return txn.env.LoadTable(e.A)
}

func (e *ConsistentTable) GetTableB(txn *Transaction) *btree.Tree {
	return txn.env.LoadTable(e.B)
}

package static

import (
	"github.com/fd/static/errors"
	"github.com/fd/static/store"
	"sync"
)

type Tx struct {
	collections map[string]*C
	terminators map[string]Terminator
	src         store.Store
	dst         store.Store
	err         errors.List

	mtx sync.Mutex
}

type Terminator interface {
	Waiter
	Open(tx *Tx) error
	Commit() error
}

func (tx *Tx) SrcStore() store.Store {
	return tx.src
}

func (tx *Tx) DstStore() store.Store {
	return tx.dst
}

func (tx *Tx) RegisterTerminator(name string, t Terminator) Terminator {
	tx.mtx.Lock()

	if tx.terminators == nil {
		tx.terminators = map[string]Terminator{}
	}

	if t := tx.terminators[name]; t != nil {
		tx.mtx.Unlock()
		return t
	}

	tx.terminators[name] = t

	err := t.Open(tx)
	if err != nil {
		tx.err.Add(err)
		tx.mtx.Unlock()
		return nil
	}

	tx.mtx.Unlock()
	return t
}

package static

import (
	"database/sql"
	"simplex.sh/errors"
	"simplex.sh/store"
	"simplex.sh/store/cas"
	"simplex.sh/store/router"
	"sync"
)

type Tx struct {
	collections   map[string]*C
	terminators   map[string]Terminator
	src           store.Store
	err           errors.List
	database      *sql.DB
	transaction   *sql.Tx
	cas_writer    *cas.Writer
	router_writer *router.Writer

	mtx sync.Mutex
}

type Terminator interface {
	Waiter
	Open(tx *Tx) error
	Commit() error
}

func (tx *Tx) SqlTx() *sql.Tx {
	return tx.transaction
}

func (tx *Tx) SrcStore() store.Store {
	return tx.src
}

func (t *Tx) Cas() *cas.Writer {
	return t.cas_writer
}

func (t *Tx) Router() *router.Writer {
	return t.router_writer
}

func (tx *Tx) RegisterTerminator(name string, t Terminator) Terminator {
	tx.mtx.Lock()
	defer tx.mtx.Unlock()

	if tx.terminators == nil {
		tx.terminators = map[string]Terminator{}
	}

	if t := tx.terminators[name]; t != nil {
		return t
	}

	tx.terminators[name] = t

	err := t.Open(tx)
	if err != nil {
		tx.err.Add(err)
		return nil
	}

	return t
}

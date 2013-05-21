package shttp

import (
	"encoding/json"
	"simplex.sh/errors"
	"simplex.sh/future"
	"simplex.sh/static"
	"sync"
)

type terminator struct {
	future.Deferred
	tx          *static.Tx
	collections []*static.C
	route_table *route_table_writer
	mtx         sync.Mutex
	wg          sync.WaitGroup
	err         errors.List
}

func terminator_for_tx(tx *static.Tx) *terminator {
	return tx.RegisterTerminator("router", &terminator{}).(*terminator)
}

func (r *terminator) Open(tx *static.Tx) error {
	r.tx = tx
	r.route_table = &route_table_writer{}
	return nil
}

func (r *terminator) Commit() error {
	r.Do(func() error {
		err := static.WaitForAll(r.collections)
		if err != nil {
			return err
		}

		r.write_route_table()
		return r.err.Normalize()
	})

	return nil
}

func (t *terminator) write_route_table() {
	w, err := t.tx.DstStore().SetBlob("route_table.json")
	if err != nil {
		t.err.Add(err)
		return
	}

	defer w.Close()

	err = json.NewEncoder(w).Encode(t.route_table)
	if err != nil {
		t.err.Add(err)
		return
	}
}

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

		r.wg.Add(1)
		go r.write_route_table()

		doc_c := make(chan *document, 10)
		r.wg.Add(1)
		go r.document_sender(doc_c)

		for i := 50; i > 0; i-- {
			r.wg.Add(1)
			go r.document_writer(doc_c)
		}

		r.wg.Wait()

		return r.err.Normalize()
	})

	return nil
}

func (t *terminator) write_route_table() {
	defer t.wg.Done()

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

func (t *terminator) document_sender(c chan *document) {
	defer t.wg.Done()
	defer close(c)

	for _, coll := range t.collections {
		l, err := coll.Len()
		if err != nil {
			t.err.Add(err)
			continue
		}

		for i := 0; i < l; i++ {
			doc_i, err := coll.At(i)
			if err != nil {
				t.err.Add(err)
				continue
			}

			doc, ok := doc_i.(*document)
			if !ok {
				continue
			}

			c <- doc
		}
	}
}

func (t *terminator) document_writer(c <-chan *document) {
	defer t.wg.Done()
	store := t.tx.DstStore()

	for doc := range c {
		w, err := store.SetBlob("blobs/" + doc.Digest)
		if err != nil {
			t.err.Add(err)
			continue
		}

		_, err = w.Write(doc.Body)
		if err != nil {
			t.err.Add(err)
			continue
		}

		err = w.Close()
		if err != nil {
			t.err.Add(err)
		}
	}
}

package cas

import (
	"bytes"
	"crypto/sha1"
	"database/sql"
	"hash"
	"sync"
)

type Writer struct {
	tx   *sql.Tx
	stmt *sql.Stmt
	keys map[string]bool
	mtx  sync.RWMutex
	err  error
}

func OpenWriter(tx *sql.Tx) (*Writer, error) {
	w := &Writer{
		tx:   tx,
		keys: make(map[string]bool, 1000),
	}

	// w.mtx.Lock()
	// defer w.mtx.Unlock()

	err := w.update_schema()
	if err != nil {
		return nil, err
	}

	{
		stmt, err := w.tx.Prepare(`INSERT INTO cas_objects (address, content) VALUES ($1, $2);`)
		if err != nil {
			return nil, err
		}

		w.stmt = stmt
	}

	w.load_known_keys()

	return w, nil
}

func (w *Writer) load_known_keys() {
	w.mtx.Lock()
	defer w.mtx.Unlock()

	rows, err := w.tx.Query(`SELECT address FROM cas_objects;`)
	if err != nil {
		w.err = err
		return
	}

	defer rows.Close()

	for rows.Next() {
		var (
			addr []byte
		)

		err := rows.Scan(&addr)
		if err != nil {
			w.err = err
			return
		}

		w.keys[string(addr)] = false
	}

	if err := rows.Err(); err != nil {
		w.err = err
		return
	}
}

type BlobWriter struct {
	w    *Writer
	buf  bytes.Buffer
	sum  hash.Hash
	addr Addr
}

func (tx *Writer) Open() *BlobWriter {
	return &BlobWriter{w: tx, sum: sha1.New()}
}

func (b *BlobWriter) Address() Addr {
	return b.addr
}

func (b *BlobWriter) Len() int {
	return b.buf.Len()
}

func (b *BlobWriter) Write(p []byte) (int, error) {
	b.sum.Write(p)
	return b.buf.Write(p)
}

func (b *BlobWriter) Close() error {
	b.w.mtx.RLock()
	defer b.w.mtx.RUnlock()

	if b.w.err != nil {
		return b.w.err
	}

	var (
		err error
	)

	b.addr = Addr(b.sum.Sum(nil))
	addr_str := string(b.addr)
	addr_map := b.w.keys

	if _, p := addr_map[addr_str]; p {
		addr_map[addr_str] = true
		return nil
	}
	addr_map[addr_str] = true

	_, err = b.w.stmt.Exec(
		[]byte(b.addr),
		b.buf.Bytes(),
	)
	if err != nil {
		return err
	}

	return nil
}

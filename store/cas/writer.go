package cas

import (
	"bytes"
	"crypto/sha1"
	"database/sql"
	"hash"
	"io"
)

type Addr []byte

type Writer interface {
	io.WriteCloser
	Address() Addr
	Len() int
}

type writer_t struct {
	transaction *sql.Tx
	buf         bytes.Buffer
	sum         hash.Hash
	addr        Addr
}

func OpenWriter(tx *sql.Tx) Writer {
	return &writer_t{transaction: tx, sum: sha1.New()}
}

func (w *writer_t) Address() Addr {
	return w.addr
}

func (w *writer_t) Len() int {
	return w.buf.Len()
}

func (w *writer_t) Write(p []byte) (int, error) {
	w.sum.Write(p)
	return w.buf.Write(p)
}

func (w *writer_t) Close() error {
	var (
		err error
	)

	w.addr = Addr(w.sum.Sum(nil))

	_, err = w.transaction.Exec(
		`INSERT INTO cas_objects (address, content) VALUES ($1, $2);`,
		[]byte(w.addr),
		w.buf.Bytes(),
	)
	if err != nil {
		return err
	}

	return nil
}

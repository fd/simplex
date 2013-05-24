package cas

import (
	"bytes"
	"database/sql"
	"io"
)

type Reader interface {
	io.ReadCloser
	Address() Addr
	Len() int
}

func OpenReader(db *sql.DB, addr Addr) (Reader, error) {
	var (
		content []byte
	)

	err := db.QueryRow(
		`SELECT content FROM cas_objects WHERE address = $1 LIMIT 1;`,
		[]byte(addr),
	).Scan(&content)
	if err != nil {
		return nil, err
	}

	return &reader_t{bytes.NewReader(content), addr}, nil
}

type reader_t struct {
	buf  *bytes.Reader
	addr Addr
}

func (w *reader_t) Address() Addr {
	return w.addr
}

func (w *reader_t) Len() int {
	return w.buf.Len()
}

func (w *reader_t) Read(b []byte) (int, error) {
	return w.buf.Read(b)
}

func (r *reader_t) Close() error {
	return nil
}

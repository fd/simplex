package cas

import (
	"bytes"
	"database/sql"
	"io"
)

type Getter interface {
	Get(Addr) (Reader, error)
}

type Reader interface {
	io.ReadCloser
	Address() Addr
	Len() int
}

type getter_t struct {
	db   *sql.DB
	stmt *sql.Stmt
}

type blob_reader_t struct {
	buf  *bytes.Reader
	addr Addr
}

func OpenGetter(db *sql.DB) (Getter, error) {
	g := &getter_t{}

	err := g.init(db)
	if err != nil {
		return nil, err
	}

	err = update_schema(db)
	if err != nil {
		return nil, err
	}

	return g, nil
}

func (g *getter_t) init(db *sql.DB) error {
	g.db = db

	{
		stmt, err := db.Prepare(`SELECT content FROM cas_objects WHERE address = $1 LIMIT 1;`)
		if err != nil {
			return err
		}

		g.stmt = stmt
	}

	return nil
}

func (g *getter_t) Get(addr Addr) (Reader, error) {
	var (
		content []byte
	)

	err := g.stmt.QueryRow([]byte(addr)).Scan(&content)
	if err != nil {
		return nil, err
	}

	return &blob_reader_t{bytes.NewReader(content), addr}, nil
}
func (w *blob_reader_t) Address() Addr {
	return w.addr
}

func (w *blob_reader_t) Len() int {
	return w.buf.Len()
}

func (w *blob_reader_t) Read(b []byte) (int, error) {
	return w.buf.Read(b)
}

func (r *blob_reader_t) Close() error {
	return nil
}

package cas

import (
	"bytes"
	"database/sql"
	"io"
	"sync"
)

type Setter interface {
	Set() Writer
}

type Writer interface {
	io.Writer
	Commit(addr Addr) error
	Abort() error
	Len() int
}

type setter_t struct {
	db   *sql.DB
	stmt *sql.Stmt

	keys map[string]bool
	mtx  sync.RWMutex
	err  error
}

type blob_writer_t struct {
	setter *setter_t
	buf    bytes.Buffer
	addr   Addr
}

func OpenSetter(db *sql.DB) (Setter, error) {
	s := &setter_t{}

	err := s.init(db)
	if err != nil {
		return nil, err
	}

	err = update_schema(db)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *setter_t) init(db *sql.DB) error {
	s.db = db
	s.keys = make(map[string]bool, 1000)

	// s.mtx.Lock()
	// defer s.mtx.Unlock()

	{
		stmt, err := s.db.Prepare(`INSERT INTO cas_objects (address, content) VALUES ($1, $2);`)
		if err != nil {
			return err
		}

		s.stmt = stmt
	}

	s.load_known_keys()

	return nil
}

func (s *setter_t) load_known_keys() {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	rows, err := s.db.Query(`SELECT address FROM cas_objects;`)
	if err != nil {
		s.err = err
		return
	}

	defer rows.Close()

	for rows.Next() {
		var (
			addr []byte
		)

		err := rows.Scan(&addr)
		if err != nil {
			s.err = err
			return
		}

		s.keys[string(addr)] = false
	}

	if err := rows.Err(); err != nil {
		s.err = err
		return
	}
}

func (s *setter_t) Set() Writer {
	return &blob_writer_t{setter: s}
}

func (b *blob_writer_t) Address() Addr {
	return b.addr
}

func (b *blob_writer_t) Len() int {
	return b.buf.Len()
}

func (b *blob_writer_t) Write(p []byte) (int, error) {
	return b.buf.Write(p)
}

func (b *blob_writer_t) Abort() error {
	return nil
}

func (b *blob_writer_t) Commit(addr Addr) error {
	b.setter.mtx.RLock()
	defer b.setter.mtx.RUnlock()

	if b.setter.err != nil {
		return b.setter.err
	}

	var (
		err error
	)

	b.addr = addr
	addr_str := string(b.addr)
	addr_map := b.setter.keys

	if _, p := addr_map[addr_str]; p {
		addr_map[addr_str] = true
		return nil
	}
	addr_map[addr_str] = true

	_, err = b.setter.stmt.Exec(
		[]byte(b.addr),
		b.buf.Bytes(),
	)
	if err != nil {
		return err
	}

	return nil
}

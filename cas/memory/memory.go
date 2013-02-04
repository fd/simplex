package memory

import (
	"bytes"
	"github.com/fd/simplex/cas"
	"io"
	"io/ioutil"
)

type store struct {
	entries map[string][]byte
}

type write_commiter struct {
	s *store
	b *bytes.Buffer
}

func New() cas.Store {
	return &store{
		entries: map[string][]byte{},
	}
}

func (s *store) Get(addr cas.Addr) (io.ReadCloser, error) {
	addr_str := addr.String()

	if d, p := s.entries[addr_str]; p {
		return ioutil.NopCloser(bytes.NewReader(d)), nil
	}

	return nil, cas.NotFound{addr}
}

func (s *store) Set() (cas.WriteCommiter, error) {
	return &write_commiter{
		s: s,
		b: bytes.NewBuffer(nil),
	}, nil
}

func (s *store) Close() error {
	s.entries = nil
	return nil
}

func (w *write_commiter) Write(p []byte) (n int, err error) {
	return w.b.Write(p)
}

func (w *write_commiter) Commit(addr cas.Addr) error {
	w.s.entries[addr.String()] = w.b.Bytes()
	return nil
}

func (w *write_commiter) Rollback() error {
	return nil
}

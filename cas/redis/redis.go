package redis

import (
	"bytes"
	"github.com/simonz05/godis/redis"
	"io"
	"io/ioutil"
	"simplex.sh/cas"
)

type store struct {
	c *redis.Client
}

type write_commiter struct {
	s *store
	b *bytes.Buffer
}

func New(addr string, db int, pwd string) cas.Store {
	return &store{
		c: redis.New(addr, db, pwd),
	}
}

func (s *store) Get(addr cas.Addr) (io.ReadCloser, error) {
	addr_str := addr.String()

	elem, err := s.c.Get(addr_str)
	if err != nil {
		return nil, err
	}

	if len(elem) == 0 {
		return nil, cas.NotFound{addr}
	}

	return ioutil.NopCloser(bytes.NewReader(elem)), nil
}

func (s *store) Set() (cas.WriteCommiter, error) {
	return &write_commiter{
		s: s,
		b: bytes.NewBuffer(nil),
	}, nil
}

func (s *store) Close() error {
	c := s.c
	s.c = nil
	return c.Quit()
}

func (w *write_commiter) Write(p []byte) (n int, err error) {
	return w.b.Write(p)
}

func (w *write_commiter) Commit(addr cas.Addr) error {
	addr_str := addr.String()

	return w.s.c.Set(addr_str, w.b.Bytes())
}

func (w *write_commiter) Rollback() error {
	return nil
}

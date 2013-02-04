package file_system

import (
	"github.com/fd/simplex/cas"
	"io"
	"io/ioutil"
	"os"
	path "path/filepath"
)

type store struct {
	root string
}

type write_commiter struct {
	s *store
	f *os.File
}

func New(root string) (cas.Store, error) {
	root, err := path.Abs(root)
	if err != nil {
		return nil, err
	}

	return &store{root}, nil
}

func (s *store) Close() error {
	return nil
}

func (s *store) Get(addr cas.Addr) (io.ReadCloser, error) {
	addr_path := s.path_for_addr(addr)

	f, err := os.Open(addr_path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, cas.NotFound{addr}
		}
		return nil, err
	}

	return f, nil
}

func (s *store) Set() (cas.WriteCommiter, error) {
	f, err := ioutil.TempFile("", "sx-fs-")
	if err != nil {
		return nil, err
	}

	return &write_commiter{s, f}, nil
}

func (w *write_commiter) Write(p []byte) (n int, err error) {
	return w.f.Write(p)
}

func (w *write_commiter) Commit(addr cas.Addr) (err error) {
	defer func() {
		if err != nil {
			w.Rollback()
		}
	}()

	err = w.f.Close()
	if err != nil {
		return err
	}

	p := w.s.path_for_addr(addr)

	// already exists?
	if _, err := os.Stat(p); err == nil {
		w.Rollback()
		return nil
	}

	// make dirs
	err = os.MkdirAll(path.Dir(p), 0755)
	if err != nil {
		return err
	}

	// move file
	tmp_name := w.f.Name()
	err = os.Rename(tmp_name, p)
	if err != nil {
		return err
	}

	return nil
}

func (w *write_commiter) Rollback() error {
	w.f.Close()
	os.Remove(w.f.Name())
	return nil
}

func (s *store) path_for_addr(addr cas.Addr) string {
	hex := addr.String()
	a := hex[0:4]
	b := hex[4:8]
	c := hex[8:]
	return path.Join(s.root, "objects", a, b, c)
}

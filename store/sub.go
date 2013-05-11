package store

import (
	"io"
	"path"
)

type sub_store_t struct {
	prefix string
	store  Store
}

func SubStore(store Store, prefix string) Store {
	return &sub_store_t{prefix, store}
}

func (s *sub_store_t) Get(name string) (io.ReadCloser, error) {
	return s.store.Get(path.Join(s.prefix, name))
}

func (s *sub_store_t) Set(name string) (io.WriteCloser, error) {
	return s.store.Set(path.Join(s.prefix, name))
}

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

func (s *sub_store_t) GetBlob(name string) (io.ReadCloser, error) {
	return s.store.GetBlob(path.Join(s.prefix, name))
}

func (s *sub_store_t) SetBlob(name string) (io.WriteCloser, error) {
	return s.store.SetBlob(path.Join(s.prefix, name))
}

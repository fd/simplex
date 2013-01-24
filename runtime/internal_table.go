package runtime

import (
	s "github.com/fd/simplex/data/storage"
	t "github.com/fd/simplex/data/trie/storage"
)

type InternalTable struct {
	Sha  s.SHA
	Name string

	txn  *Transaction
	trie t.T
}

func (t *InternalTable) Get(key, value interface{}) bool {
	return false
}

func (t *InternalTable) Set(key, value interface{}) (prev, curr s.SHA, changed bool) {
	return s.ZeroSHA, s.ZeroSHA, false
}

func (t *InternalTable) Del(key interface{}) (prev s.SHA, changed bool) {
	return s.ZeroSHA, false
}

func (t *InternalTable) Commit() (prev, curr s.SHA, changed bool) {
	t.Sha = t.trie.Commit()
	return t.txn.tables.Set(t.Name, t)
}

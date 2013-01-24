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

func (it *InternalTable) setup() {
	if it.Sha != s.ZeroSHA {
		trie, ok := t.Get(it.txn.env.store, it.Sha)
		if !ok {
			panic("corrupted data store")
		}
		it.trie = trie
	}

	if it.trie == nil {
		it.trie = t.New(it.txn.env.store)
	}
}

func (t *InternalTable) Get(key, value interface{}) bool {
	bin_key := consistent_rep(key)

	value_sha, found := t.trie.Lookup(bin_key)
	if !found {
		return false
	}

	return t.txn.env.store.Get(value_sha, value)
}

func (t *InternalTable) Set(key, value interface{}) (prev, curr s.SHA, changed bool) {
	bin_key := consistent_rep(key)

	curr = t.txn.env.store.Set(value)

	prev, inserted := t.trie.Insert(bin_key, curr)
	if !inserted {
		panic("insert failed")
	}

	return prev, curr, (curr != prev)
}

func (t *InternalTable) Del(key interface{}) (prev s.SHA, changed bool) {
	panic("not implemented")
}

func (t *InternalTable) Commit() (prev, curr s.SHA, changed bool) {
	t.Sha = t.trie.Commit()

	// the _tables table is special
	if t.Name == "_tables" {
		return
	}

	return t.txn.Tables.Set(t.Name, t)
}

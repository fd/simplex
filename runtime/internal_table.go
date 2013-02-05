package runtime

import (
	"bytes"
	"github.com/fd/simplex/cas"
	"github.com/fd/simplex/cas/btree"
)

type InternalTable struct {
	Addr cas.Addr
	Name string

	env  *Environment
	txn  *Transaction
	tree *btree.Tree
}

type Iterator interface {
	Next() (sha s.SHA, done bool)
}

type KeyValue struct {
	KeyCompare []byte
	ValueSha   s.SHA
}

func (it *InternalTable) setup() {
	if it.txn != nil {
		it.env = it.txn.env
	}

	if it.Sha != s.ZeroSHA {
		trie, ok := t.Get(it.env.store, it.Sha)
		if !ok {
			panic("corrupted data store")
		}
		it.trie = trie
	}

	if it.trie == nil {
		it.trie = t.New(it.env.store)
	}
}

func (t *InternalTable) Iter() Iterator {
	return t.trie.Iter()
}

func (t *InternalTable) GetKeyValueSHA(key []byte) (s.SHA, bool) {
	return t.trie.Lookup(key)
}

func (t *InternalTable) Get(key, value interface{}) bool {
	bin_key := consistent_rep(key)

	kv_sha, found := t.GetKeyValueSHA(bin_key)
	if !found {
		return false
	}

	kv := KeyValue{}
	if !t.env.store.Get(kv_sha, &kv) {
		return false
	}

	return t.env.store.Get(kv.ValueSha, value)
}

func (t *InternalTable) Set(kv *KeyValue) (prev, curr s.SHA, changed bool) {
	if t.txn == nil {
		panic("read only trie")
	}

	curr = t.env.store.Set(kv)

	prev, inserted := t.trie.Insert(kv.KeyCompare, curr)
	if !inserted {
		panic("insert failed")
	}

	curr_bytes := [20]byte(curr)
	prev_bytes := [20]byte(prev)

	return prev, curr, bytes.Compare(curr_bytes[:], prev_bytes[:]) != 0
}

func (t *InternalTable) Del(kv *KeyValue) (prev s.SHA, changed bool) {
	if t.txn == nil {
		panic("read only trie")
	}

	return t.trie.Remove(kv.KeyCompare)
}

func (t *InternalTable) Commit() (prev, curr s.SHA, changed bool) {
	if t.txn == nil {
		panic("read only trie")
	}

	t.Sha = t.trie.Commit()

	// the _tables table is special
	if t.Name == "_tables" {
		return
	}

	return t.txn.Tables.Set(&KeyValue{
		consistent_rep(t.Name),
		t.txn.env.store.Set(t),
	})
}

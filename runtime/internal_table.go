package runtime

import (
	"bytes"
	s "github.com/fd/simplex/data/storage"
	t "github.com/fd/simplex/data/trie/storage"
)

type InternalTable struct {
	Sha  s.SHA
	Name string

	txn  *Transaction
	trie t.T
}

type KeyValue struct {
	KeyCompare []byte
	ValueSha   s.SHA
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

	kv_sha, found := t.trie.Lookup(bin_key)
	if !found {
		return false
	}

	kv := KeyValue{}
	if !t.txn.env.store.Get(kv_sha, &kv) {
		return false
	}

	return t.txn.env.store.Get(kv.ValueSha, value)
}

func (t *InternalTable) Set(kv *KeyValue) (prev, curr s.SHA, changed bool) {
	curr = t.txn.env.store.Set(kv)

	prev, inserted := t.trie.Insert(kv.KeyCompare, curr)
	if !inserted {
		panic("insert failed")
	}

	curr_bytes := [20]byte(curr)
	prev_bytes := [20]byte(prev)

	return prev, curr, bytes.Compare(curr_bytes[:], prev_bytes[:]) != 0
}

func (t *InternalTable) Del(kv *KeyValue) (prev s.SHA, changed bool) {
	return t.trie.Remove(kv.KeyCompare)
}

func (t *InternalTable) Commit() (prev, curr s.SHA, changed bool) {
	t.Sha = t.trie.Commit()

	// the _tables table is special
	if t.Name == "_tables" {
		return
	}

	return t.txn.Tables.Set(&KeyValue{
		[]byte(t.Name),
		t.txn.env.store.Set(&t),
	})
}

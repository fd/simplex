package data

import (
	"fmt"
	"github.com/fd/w/data/trie"
)

type SHA [20]byte

type Emiter interface {
	Emit(key, value Value)
}

type KeyValue struct {
	KeySHA   SHA
	ValueSHA SHA
}

type (
	ScalarMapFunc    func(e Emiter, key, value Value)
	ScalarReduceFunc func(e Emiter, key Value, values []Value)
	ScalarMergeFunc  func(e Emiter, key Value, left, right []Value)
)

type MapReduce struct {
	Id     SHA
	Map    ScalarMapFunc
	Reduce ScalarReduceFunc
}

type MapReduceTransaction struct {
	// the source collection
	SourceSHA SHA

	// the previously calculated collection
	PreviousSHA SHA

	// SHA(SHA(key) + SHA(value)) of added members (in the source)
	Added []SHA

	// keys of removed members (in the source)
	Removed []SHA
}

func (mr *MapReduce) Run(txn MapReduceTransaction, store Store) (added, removed []SHA) {
	var (
		state *MapReduceState
		found bool
	)

	found = store.Get(txn.PreviousSHA, &state)
	if !found {
		state = &MapReduceState{
			MapStage:    trie.New(),
			ReduceStage: trie.New(),
		}
		txn.Removed = nil

		// load all source kvs
		var source Collection

		found = store.Get(txn.SourceSHA, &source)
		if !found {
			panic(fmt.Sprintf("MapReduce requires a source collection (not found: %s)", txn.SourceSHA))
		}

		txn.Added = source.MemberSHAs()
	}

	state.SHA = HashValue([]SHA{txn.SourceSHA, mr.Id})

	var (
		need_reduce   map[string]bool
		remove_reduce map[string]bool
	)

	for _, kv_sha := range txn.Added {
		var (
			kv          *KeyValue
			key         Value
			val         Value
			map_key_str []byte

			found  bool
			emiter = &map_emiter{}
		)

		found = store.Get(kv_sha, &kv)
		if !found {
			panic(fmt.Sprintf("corrupted datastore: missing KeyValue(%s)", kv_sha))
		}

		found = store.Get(kv.KeySHA, &key)
		if !found {
			panic(fmt.Sprintf("corrupted datastore: missing Key(%s)", kv.KeySHA))
		}

		found = store.Get(kv.ValueSHA, &val)
		if !found {
			panic(fmt.Sprintf("corrupted datastore: missing Value(%s)", kv.ValueSHA))
		}

		map_key_str = CompairBytes(key)
		mr.Map(emiter, key, val)

		for _, pair := range emiter.pairs {
			kv := &KeyValue{
				KeySHA:   store.Set(pair.Key),
				ValueSHA: store.Set(pair.Val),
			}

			kv_sha := store.Set(kv)
			reduce_key_str := CompairBytes(pair.Key)

			reduce_bucket_i, found := state.MapStage.Lookup(reduce_key_str)
			reduce_bucket, ok := reduce_bucket_i.(*trie.T)
			if !ok {
				panic(fmt.Sprintf("corrupted datastore: Invalid reduce bucket (%v)", pair.Key))
			}
			if !found {
				reduce_bucket = trie.New()
				state.MapStage.Insert(reduce_key_str, reduce_bucket)
			}
			reduce_bucket.Insert(map_key_str, kv_sha)

			need_reduce[string(reduce_key_str)] = false // false == partial
		}
	}

	for _, kv_sha := range txn.Removed {
		var (
			kv          *KeyValue
			key         Value
			val         Value
			map_key_str []byte

			found  bool
			emiter = &map_emiter{}
		)

		found = store.Get(kv_sha, &kv)
		if !found {
			panic(fmt.Sprintf("corrupted datastore: missing KeyValue(%s)", kv_sha))
		}

		found = store.Get(kv.KeySHA, &key)
		if !found {
			panic(fmt.Sprintf("corrupted datastore: missing Key(%s)", kv.KeySHA))
		}

		found = store.Get(kv.ValueSHA, &val)
		if !found {
			panic(fmt.Sprintf("corrupted datastore: missing Value(%s)", kv.ValueSHA))
		}

		map_key_str = CompairBytes(key)
		mr.Map(emiter, key, val)

		for _, pair := range emiter.pairs {
			reduce_key_str := CompairBytes(pair.Key)

			reduce_bucket_i, found := state.MapStage.Lookup(reduce_key_str)
			if !found {
				remove_reduce[string(reduce_key_str)] = true
				continue // ignore
			}

			reduce_bucket, ok := reduce_bucket_i.(*trie.T)
			if !ok {
				panic(fmt.Sprintf("corrupted datastore: Invalid reduce bucket (%v)", pair.Key))
			}

			reduce_bucket.Remove(map_key_str)

			if reduce_bucket.Len() == 0 {
				state.MapStage.Remove(reduce_key_str)
				remove_reduce[string(reduce_key_str)] = true
				continue
			}

			need_reduce[string(reduce_key_str)] = true // true == full
		}
	}

	for reduce_key_str := range remove_reduce {
		delete(need_reduce, reduce_key_str)
		old_kv_sha_i, f := state.ReduceStage.Remove([]byte(reduce_key_str))
		old_kv_sha, _ := old_kv_sha_i.(SHA)
		if f {
			removed = append(removed, old_kv_sha)
		}
	}

	// for reduce_key_str, _full := range need_reduce {
	for reduce_key_str := range need_reduce {
		var (
			key        Value
			val_sha_is []interface{}
			vals       []Value
			emiter     = &reduce_emiter{}

			reduce_key_bytes = []byte(reduce_key_str)
			key_sha          = HashCompairBytes(reduce_key_bytes)
		)

		found = store.Get(key_sha, &key)
		if !found {
			panic(fmt.Sprintf("corrupted datastore: missing Key(%s)", key_sha))
		}

		reduce_bucket_i, found := state.MapStage.Lookup(reduce_key_bytes)
		reduce_bucket, ok := reduce_bucket_i.(*trie.T)
		if !ok {
			panic(fmt.Sprintf("corrupted datastore: Invalid reduce bucket (%v)", key))
		}
		if !found {
			panic(fmt.Sprintf("corrupted datastore: Missing reduce bucket (%v)", key))
		}

		val_sha_is = reduce_bucket.Values()
		vals = make([]Value, 0, len(val_sha_is))
		for _, val_sha_i := range val_sha_is {
			var val Value

			val_sha, ok := val_sha_i.(SHA)
			if !ok {
				panic(fmt.Sprintf("corrupted datastore: Invalid reduce bucket (%v)", key))
			}

			found = store.Get(val_sha, &val)
			if !found {
				panic(fmt.Sprintf("corrupted datastore: missing Value(%s)", val_sha))
			}

			vals = append(vals, val)
		}

		emiter.key = key
		mr.Reduce(emiter, key, vals)

		kv := &KeyValue{
			KeySHA:   store.Set(key),
			ValueSHA: store.Set(emiter.val),
		}
		kv_sha := store.Set(kv)

		old_kv_sha_i, f := state.ReduceStage.Insert(reduce_key_bytes, kv_sha)
		old_kv_sha, _ := old_kv_sha_i.(SHA)
		if old_kv_sha != kv_sha {
			if f {
				removed = append(removed, old_kv_sha)
			}
			added = append(added, kv_sha)
		}
	}

	// continue with added and removed
	return added, removed
}

type map_emiter struct {
	pairs []map_emiter_pair
}

type map_emiter_pair struct {
	Key Value
	Val Value
}

func (e *map_emiter) Emit(key, val Value) {
	e.pairs = append(e.pairs, map_emiter_pair{key, val})
}

type reduce_emiter struct {
	key Value
	val Value
	set bool
}

func (e *reduce_emiter) Emit(key, val Value) {
	if e.set {
		panic(fmt.Sprintf("reduce: Only one key value pair can be emited during a reduce."))
	}
	if key != e.key {
		panic(fmt.Sprintf("reduce: Emited key must match the input key. (%v != %v)", key, e.key))
	}
	e.val = val
	e.set = true
}

type MapReduceState struct {
	SHA         SHA
	MapStage    *trie.T
	ReduceStage *trie.T
}

type Collection interface {
	MemberSHAs() []SHA
}

type Store interface {
	Get(sha SHA, value interface{}) bool
	Set(value interface{}) SHA
}

func HashValue(v interface{}) SHA {
	return HashCompairBytes(CompairBytes(v))
}

func HashCompairBytes(v []byte) SHA {
	return SHA{}
}

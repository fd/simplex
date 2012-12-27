package transformations

import (
	"fmt"
	"github.com/fd/w/data/ident"
	"github.com/fd/w/data/storage"
	"github.com/fd/w/data/transaction"
	trie "github.com/fd/w/data/trie/storage"
	"github.com/fd/w/data/value"
)

type Emiter interface {
	Emit(key, value value.Any)
}

type KeyValue struct {
	KeySHA   ident.SHA
	ValueSHA ident.SHA
}

type MapReduce struct {
	Id     ident.SHA
	Map    ScalarMapFunc
	Reduce ScalarReduceFunc
}

func load_kv(store *storage.S, kv_sha ident.SHA) (*KeyValue, value.Any, value.Any) {
	var (
		kv  *KeyValue
		key value.Any
		val value.Any
	)

	if !store.Get(kv_sha, &kv) {
		panic(fmt.Sprintf("corrupted datastore: missing KeyValue(%s)", kv_sha))
	}

	if !store.Get(kv.KeySHA, &key) {
		panic(fmt.Sprintf("corrupted datastore: missing Key(%s)", kv.KeySHA))
	}

	if !store.Get(kv.ValueSHA, &val) {
		panic(fmt.Sprintf("corrupted datastore: missing Value(%s)", kv.ValueSHA))
	}

	return kv, key, val
}

type bucket_t struct {
	Key     ident.SHA
	Value   ident.SHA
	Mapped  ident.SHA
	Reduced ident.SHA

	key     value.Any
	mapped  trie.T
	reduced trie.T
}

func find_bucket(mr_trie trie.T, cache map[string]*bucket_t, key value.Any) *bucket_t {
	store := mr_trie.Store()
	key_str := ident.CompairBytes(key)

	bucket, found := cache[string(key_str)]
	if found {
		return bucket
	}

	bucket_sha, found := mr_trie.Lookup(key_str)
	if !found {
		bucket = &bucket_t{
			key:     key,
			mapped:  trie.New(store),
			reduced: trie.New(store),
		}
		cache[string(key_str)] = bucket
		return bucket
	}

	if !store.Get(bucket_sha, &bucket) {
		panic(fmt.Sprintf("corrupted datastore: Invalid map_reduce bucket (%s)", bucket_sha))
	}

	mapped, found := trie.Get(store, bucket.Mapped)
	if !found {
		panic(fmt.Sprintf("corrupted datastore: Invalid map_reduce bucket (%s)", bucket_sha))
	}
	bucket.mapped = mapped

	reduced, found := trie.Get(store, bucket.Reduced)
	if !found {
		panic(fmt.Sprintf("corrupted datastore: Invalid map_reduce bucket (%s)", bucket_sha))
	}
	bucket.reduced = reduced

	bucket.key = key

	cache[string(key_str)] = bucket
	return bucket
}

type map_op struct {
	typ    int
	kv_sha ident.SHA
}

type map_emit_op struct {
	typ   int
	m_key []byte
	r_key value.Any
	val   value.Any
}

type sort_emit_op struct {
	r_key  []byte
	bucket *bucket_t
}

type map_emiter struct {
	typ      int
	m_key    []byte
	out_chan chan<- map_emit_op
}

func (e map_emiter) Emit(key, val value.Any) {
	e.out_chan <- map_emit_op{e.typ, e.m_key, key, val}
}

func mapper(mapf ScalarMapFunc, store *storage.S, in_chan <-chan map_op, out_chan chan<- map_emit_op) {
	defer close(out_chan)
	for op := range in_chan {
		kv, key, val := load_kv(store, op.kv_sha)
		mapf(map_emiter{op.typ, ident.CompairBytes(key), out_chan}, key, val)
	}
}

func sorter(mr_trie trie.T, in_chan <-chan map_emit_op, out_chan chan<- sort_emit_op) {
	defer close(out_chan)

	var (
		cache = map[string]*bucket_t{}
		store = mr_trie.Store()
	)

	for emit := range in_chan {

		bucket := find_bucket(
			mr_trie,
			cache,
			emit.r_key,
		)

		if emit.typ == 1 { // insert
			bucket.mapped.Insert(
				emit.m_key,
				store.Set(emit.val),
			)
		}

		if emit.typ == 2 { // remove
			bucket.mapped.Remove(emit.m_key)
		}

	}

	for key_str, bucket := range cache {

		// store mapped trie
		bucket.Mapped = bucket.mapped.Commit()

		// send to reducer
		out_chan <- sort_emit_op{[]byte(key_str), bucket}

	}
}

func cleaner(mr_trie trie.T, in_chan <-chan sort_emit_op, out_chan chan<- sort_emit_op, removed chan<- ident.SHA) {
	defer close(out_chan)

	for op := range in_chan {
		if op.bucket.mapped.Empty() {
			mr_trie.Remove(op.r_key)
			removed <- op.r_key
		} else {
			out_chan <- op
		}
	}
}

func load_state(mr *MapReduce, store *storage.S, action transaction.Action) trie.T {
	var (
		state trie.T
	)

	state, found := trie.Get(store, action.PreviousSHA)
	if found {
		return state
	}

	state = trie.New(store)
	action.Removed = nil

	// load all source kvs
	var source Collection

	found = store.Get(action.SourceSHA, &source)
	if !found {
		panic(fmt.Sprintf("MapReduce requires a source collection (not found: %s)", action.SourceSHA))
	}

	action.Added = source.MemberSHAs()
	return state
}

func reducer(mr_trie trie.T, reducef ScalarReduceFunc, in_chan <-chan sort_emit_op, out_chan chan<- bool) {
	defer func() {
		out_chan <- true
		close(out_chan)
	}()

	store := mr_trie.Store()

	for op := range in_chan {

		// reduce mapped values
		val_sha := op.bucket.mapped.Reduce(op.bucket.reduced, func(reductions []ident.SHA, val_sha ident.SHA) ident.SHA {

			var (
				val  value.Any
				vals = make([]value.Any, 0, len(reductions)+1)
			)

			if val_sha != ident.ZeroSHA {
				val = nil
				if !store.Get(val_sha, &val) {
					panic(fmt.Sprintf("corrupted datastore: missing value (%s)", val_sha))
				}
				vals = append(vals, val)
			}

			for _, val_sha := range reductions {
				val = nil
				if !store.Get(val_sha, &val) {
					panic(fmt.Sprintf("corrupted datastore: missing value (%s)", val_sha))
				}
				vals = append(vals, val)
			}

			emiter := &reduce_emiter{}
			reducef(emiter, op.bucket.key, vals)

			return store.Set(emiter.val)

		})

		// store reduce cache
		op.bucket.Reduced = op.bucket.reduced.Commit()
		op.bucket.Key = store.Set(op.bucket.key)
		op.bucket.Value = val_sha

		// store op.bucket
		bucket_sha := store.Set(op.bucket)
		mr_trie.Insert(op.r_key, bucket_sha)
	}
}

func push_changes(action transaction.Action, out_chan chan<- map_op) {
	defer close(out_chan)

	for _, kv_sha := range action.Added {
		out_chan <- map_op{1, kv_sha}
	}

	for _, kv_sha := range action.Removed {
		out_chan <- map_op{2, kv_sha}
	}
}

func collect_shas(in_chan <-chan ident.SHA, s *[]ident.SHA) {
	a := []ident.SHA{}

	for sha := range in_chan {
		a = append(a, sha)
	}

	*s = a
}

func (mr *MapReduce) Run(action transaction.Action, store *storage.S) []transaction.Action {
	var (
		found   bool
		added   []ident.SHA
		removed []ident.SHA

		removed_shas = make(chan ident.SHA, 10)
		added_shas   = make(chan ident.SHA, 10)

		state        = load_state(mr, store, action)
		map_chan     = make(chan map_op, 10)
		mapped_chan  = make(chan map_emit_op, 10)
		sorted_chan  = make(chan sort_emit_op, 10)
		cleaned_chan = make(chan sort_emit_op, 10)
		done_chan    = make(chan bool, 1)
	)

	{
		go push_changes(action, map_chan)

		// concurency can be increased (when defer close() is removed)
		go mapper(mr.Map, store, map_chan, mapped_chan)

		go sorter(state, mapped_chan, sorted_chan)

		go cleaner(state, sorted_chan, cleaned_chan, removed_shas)

		go reducer(state, mr.Reduce, cleaned_chan, done_chan)

		go collect_shas(added_shas, &added)

		go collect_shas(removed_shas, &removed)
	}

	<-done_chan

	// commit state.trie
	state_sha := state.Commit()

	return []transaction.Action{
		{
			SourceSHA: state_sha,
			Added:     added,
			Removed:   removed,
		},
	}
}

type reduce_emiter struct {
	key value.Any
	val value.Any
	set bool
}

func (e *reduce_emiter) Emit(key, val value.Any) {
	if e.set {
		panic(fmt.Sprintf("reduce: Only one key value pair can be emited during a reduce."))
	}
	if key != e.key {
		panic(fmt.Sprintf("reduce: Emited key must match the input key. (%v != %v)", key, e.key))
	}
	e.val = val
	e.set = true
}

type Collection interface {
	MemberSHAs() []ident.SHA
}

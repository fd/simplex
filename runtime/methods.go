package runtime

import (
	"github.com/fd/simplex/cas"
)

func DeclareTable(name string) Deferred {
	return &table_op{name}
}

/*
  type V view[]M
  V.select(func(M)bool) -> V
*/
func Select(v IndexedView, f select_func, name string) Deferred {
	return &select_op{src: v, fun: f, name: name}
}

/*
  type V view[]M
  V.reject(func(M)bool) -> V
*/
func Reject(v IndexedView, f reject_func, name string) Deferred {
	return &reject_op{src: v, fun: f, name: name}
}

/*
  type V view[]M
  V.detect(func(M)bool) -> M

  TODO(fd) detect is only valid in a transactional function.
*/
func Detect(v IndexedView, f func(interface{}) bool, name string) interface{} {
	panic("not yet implemented")
}

/*
  type V view[]M
  type W view[]N
  V.collect(func(M)N) -> W
  (Note: the key type remains unchanged)
*/
func Collect(v IndexedView, f collect_func, name string) Deferred {
	return &collect_op{src: v, fun: f, name: name}
}

/*
  type V view[]M
  V.inject(func(M, []N)N) -> N

  TODO(fd) inject is only valid in a transactional function.
*/
func Inject(v IndexedView, f func(interface{}, []interface{}) interface{}, name string) interface{} {
	panic("not yet implemented")
}

/*
  type V view[]M
  type W view[N]view[]M
  V.group(func(M)N) -> W
  (Note: the key type of the inner view remains unchanged)
*/
func Group(v IndexedView, f group_func, name string) Deferred {
	return &group_op{src: v, fun: f, name: name}
}

/*
  type V view[]M
  type W view[N]M
  V.index(func(M)N) -> W
  (Note: the member values remain unchanged)

  v.index(f) is equivalent to v.group(f).collect(func(v view[]M)M{ return v.detect(func(_){return true}) })
*/
func Index(v IndexedView, f index_func, name string) Deferred {
	return &index_op{src: v, fun: f, name: name}
}

/*
  type V view[]M
  V.sort(func(M)N) -> V
  (Note: the key type is lost)
*/
func Sort(v IndexedView, f sort_func, name string) Deferred {
	return &sort_op{src: v, fun: f, name: name}
}

func Union(v ...Deferred) Deferred {
	panic("not yet implemented")
}

type (
	table_op struct {
		name string
	}

	select_func func(*Context, SHA) bool
	select_op   struct {
		name string
		src  IndexedView
		fun  select_func
	}

	reject_func func(interface{}) bool
	reject_op   struct {
		name string
		src  IndexedView
		fun  reject_func
	}

	collect_func func(*Context, SHA) SHA
	collect_op   struct {
		name string
		src  IndexedView
		fun  collect_func
	}

	group_func func(*Context, SHA) SHA
	group_op   struct {
		name string
		src  IndexedView
		fun  group_func
	}

	index_func func(*Context, SHA) SHA
	index_op   struct {
		name string
		src  IndexedView
		fun  index_func
	}

	sort_func func(*Context, SHA) SHA
	sort_op   struct {
		name string
		src  IndexedView
		fun  sort_func
	}
)

func (op *table_op) DeferredId() string   { return op.name }
func (op *select_op) DeferredId() string  { return op.name }
func (op *reject_op) DeferredId() string  { return op.name }
func (op *collect_op) DeferredId() string { return op.name }
func (op *group_op) DeferredId() string   { return op.name }
func (op *index_op) DeferredId() string   { return op.name }
func (op *sort_op) DeferredId() string    { return op.name }

func (op *table_op) Resolve(txn *Transaction, events chan<- Event) {
	table := txn.GetTable(op.name)

	for _, change := range txn.changes {
		if change.Table != op.name {
			continue
		}

		switch change.Kind {
		case SET:
			kv := KeyValue{
				KeyCompare: consistent_rep(change.Key),
				ValueSha:   txn.env.store.Set(change.Value),
			}
			old_kv_sha, new_kv_sha, changed := table.Set(&kv)
			if changed {
				events <- &ev_CHANGE{op.name, old_kv_sha, new_kv_sha}
			}

		case UNSET:
			kv := KeyValue{
				KeyCompare: consistent_rep(change.Key),
			}
			old_kv_sha, deleted := table.Del(&kv)
			if deleted {
				events <- &ev_CHANGE{op.name, old_kv_sha, storage.ZeroSHA}
			}

		}
	}

	old_tab_sha, new_tab_sha, _ := table.Commit()
	events <- &EvConsistent{op.name, old_tab_sha, new_tab_sha}
}

func (op *select_op) Resolve(txn *Transaction, events chan<- Event) {
	var (
		src_event = txn.Resolve(op.src)
		table     = txn.GetTable(op.name)
	)

	for event := range src_event {
		i_change, ok := event.(*ev_CHANGE)
		if !ok {
			continue
		}

		var (
			o_change = &ev_CHANGE{op.name, i_change.a, i_change.b}
			kv_a     KeyValue
			kv_b     KeyValue
		)

		if !o_change.a.IsZero() {
			// lookup key/value
			if !txn.env.store.Get(o_change.a, &kv_a) {
				panic("corrupt data store")
			}

			if !op.fun(&Context{txn}, SHA(kv_a.ValueSha)) {
				o_change.a = storage.ZeroSHA
			}
		}

		if !o_change.b.IsZero() {
			// lookup key/value
			if !txn.env.store.Get(o_change.b, &kv_b) {
				panic("corrupt data store")
			}

			if !op.fun(&Context{txn}, SHA(kv_b.ValueSha)) {
				o_change.b = storage.ZeroSHA
			}
		}

		// ignore unchanged data
		if o_change.a.IsZero() && o_change.b.IsZero() {
			continue
		}

		if !o_change.a.IsZero() {
			// remove kv from table
			_, deleted := table.Del(&kv_a)
			if !deleted {
				o_change.a = storage.ZeroSHA
			}
		}

		if !o_change.b.IsZero() {
			// insert kv into table
			_, _, changed := table.Set(&kv_b)
			if !changed {
				o_change.b = storage.ZeroSHA
			}
		}

		// ignore unchanged data
		if o_change.a.IsZero() && o_change.b.IsZero() {
			continue
		}

		events <- o_change
	}

	old_tab_sha, new_tab_sha, _ := table.Commit()
	events <- &EvConsistent{op.name, old_tab_sha, new_tab_sha}
}

func (op *collect_op) Resolve(txn *Transaction, events chan<- Event) {
	var (
		src_event = txn.Resolve(op.src)
		table     = txn.GetTable(op.name)
	)

	for event := range src_event {
		i_change, ok := event.(*ev_CHANGE)
		if !ok {
			continue
		}

		// removed
		if i_change.b.IsZero() {
			var kv_a KeyValue
			// lookup key/value
			if !txn.env.store.Get(i_change.a, &kv_a) {
				panic("corrupt data store")
			}

			if kv_sha, ok := table.GetKeyValueSHA(kv_a.KeyCompare); ok {
				table.Del(&kv_a)
				events <- &ev_CHANGE{op.name, kv_sha, storage.ZeroSHA}
			}

			continue
		}

		{ // added or updated
			var kv_a, kv_b KeyValue
			// lookup key/value
			if !txn.env.store.Get(i_change.b, &kv_a) {
				panic("corrupt data store")
			}

			kv_b.KeyCompare = kv_a.KeyCompare
			kv_b.ValueSha = storage.SHA(op.fun(&Context{txn}, SHA(kv_a.ValueSha)))

			prev_kv_sha, curr_kv_sha, changed := table.Set(&kv_b)
			if changed {
				events <- &ev_CHANGE{op.name, prev_kv_sha, curr_kv_sha}
			}
		}
	}

	old_tab_sha, new_tab_sha, _ := table.Commit()
	events <- &EvConsistent{op.name, old_tab_sha, new_tab_sha}
}

func (op *reject_op) Resolve(txn *Transaction, events chan<- Event) {}
func (op *group_op) Resolve(txn *Transaction, events chan<- Event)  {}
func (op *index_op) Resolve(txn *Transaction, events chan<- Event)  {}
func (op *sort_op) Resolve(txn *Transaction, events chan<- Event)   {}

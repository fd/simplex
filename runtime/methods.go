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

	select_func func(*Context, cas.Addr) bool
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

	collect_func func(*Context, cas.Addr) cas.Addr
	collect_op   struct {
		name string
		src  IndexedView
		fun  collect_func
	}

	group_func func(*Context, cas.Addr) cas.Addr
	group_op   struct {
		name string
		src  IndexedView
		fun  group_func
	}

	index_func func(*Context, cas.Addr) cas.Addr
	index_op   struct {
		name string
		src  IndexedView
		fun  index_func
	}

	sort_func func(*Context, cas.Addr) cas.Addr
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
			var (
				key_coll      []byte
				key_addr      cas.Addr
				elt_addr      cas.Addr
				prev_elt_addr cas.Addr
				err           error
			)

			key_coll = cas.Collate(change.Key)

			key_addr, err = cas.Encode(txn.env.Store, change.Key, -1)
			if err != nil {
				panic("runtime: " + err.Error())
			}

			elt_addr, err = cas.Encode(txn.env.Store, change.Elt, -1)
			if err != nil {
				panic("runtime: " + err.Error())
			}

			prev_elt_addr, err = table.Set(key_coll, key_addr, elt_addr)
			if err != nil {
				panic("runtime: " + err.Error())
			}

			if cas.CompareAddr(prev_elt_addr, elt_addr) != 0 {
				events <- &ev_CHANGE{op.name, key_coll, key_addr, prev_elt_addr, elt_addr}
			}

		case UNSET:
			var (
				key_coll []byte
				key_addr cas.Addr
				elt_addr cas.Addr
				err      error
			)

			key_coll = cas.Collate(change.Key)

			key_addr, elt_addr, err = table.Del(key_coll)
			if err != nil {
				panic("runtime: " + err.Error())
			}

			if key_addr != nil || elt_addr != nil {
				events <- &ev_CHANGE{op.name, key_coll, key_addr, elt_addr, nil}
			}

		}
	}

	tab_addr_a, tab_addr_b := txn.CommitTable(op.name, table)
	events <- &EvConsistent{op.name, tab_addr_a, tab_addr_b}
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
			o_change = &ev_CHANGE{op.name, i_change.collated_key, i_change.key, i_change.a, i_change.b}
		)

		if o_change.a != nil {
			if !op.fun(&Context{txn}, o_change.a) {
				o_change.a = nil
			}
		}

		if o_change.b != nil {
			if !op.fun(&Context{txn}, o_change.b) {
				o_change.b = nil
			}
		}

		// ignore unchanged data
		if o_change.a == nil && o_change.b == nil {
			continue
		}

		if o_change.a != nil {
			// remove kv from table
			_, prev, err := table.Del(o_change.collated_key)
			if err != nil {
				panic("runtime: " + err.Error())
			}
			if prev != nil {
				o_change.a = nil
			}
		}

		if o_change.b != nil {
			// insert kv into table
			prev, err := table.Set(o_change.collated_key, o_change.key, o_change.b)
			if err != nil {
				panic("runtime: " + err.Error())
			}
			if cas.CompareAddr(prev, o_change.b) == 0 {
				o_change.b = nil
			}
		}

		// ignore unchanged data
		if o_change.a == nil && o_change.b == nil {
			continue
		}

		events <- o_change
	}

	tab_addr_a, tab_addr_b := txn.CommitTable(op.name, table)
	events <- &EvConsistent{op.name, tab_addr_a, tab_addr_b}
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
		if i_change.b == nil {
			prev_key_addr, prev_elt_addr, err := table.Del(i_change.collated_key)
			if err != nil {
				panic("runtime: " + err.Error())
			}

			if prev_key_addr != nil && prev_elt_addr != nil {
				events <- &ev_CHANGE{op.name, i_change.collated_key, prev_key_addr, prev_elt_addr, nil}
			}

			continue
		}

		{ // added or updated
			curr_elt_addr := op.fun(&Context{txn}, i_change.b)

			prev_elt_addr, err := table.Set(i_change.collated_key, i_change.key, curr_elt_addr)
			if err != nil {
				panic("runtime: " + err.Error())
			}
			if cas.CompareAddr(prev_elt_addr, curr_elt_addr) != 0 {
				events <- &ev_CHANGE{op.name, i_change.collated_key, i_change.key, prev_elt_addr, curr_elt_addr}
			}
		}
	}

	tab_addr_a, tab_addr_b := txn.CommitTable(op.name, table)
	events <- &EvConsistent{op.name, tab_addr_a, tab_addr_b}
}

func (op *reject_op) Resolve(txn *Transaction, events chan<- Event) {}
func (op *group_op) Resolve(txn *Transaction, events chan<- Event)  {}
func (op *index_op) Resolve(txn *Transaction, events chan<- Event)  {}
func (op *sort_op) Resolve(txn *Transaction, events chan<- Event)   {}

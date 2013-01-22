package runtime

func DeclareTable(name string) Deferred {
	return &table_op{name}
}

/*
  type V view[]M
  V.select(func(M)bool) -> V
*/
func Select(v IndexedView, f select_func) Deferred {
	return &select_op{src: v, fun: f}
}

/*
  type V view[]M
  V.reject(func(M)bool) -> V
*/
func Reject(v IndexedView, f reject_func) Deferred {
	return &reject_op{src: v, fun: f}
}

/*
  type V view[]M
  V.detect(func(M)bool) -> M

  TODO(fd) detect is only valid in a transactional function.
*/
func Detect(v IndexedView, f func(interface{}) bool) interface{} {
	panic("not yet implemented")
}

/*
  type V view[]M
  type W view[]N
  V.collect(func(M)N) -> W
  (Note: the key type remains unchanged)
*/
func Collect(v IndexedView, f collect_func) Deferred {
	return &collect_op{src: v, fun: f}
}

/*
  type V view[]M
  V.inject(func(M, []N)N) -> N

  TODO(fd) inject is only valid in a transactional function.
*/
func Inject(v IndexedView, f func(interface{}, []interface{}) interface{}) interface{} {
	panic("not yet implemented")
}

/*
  type V view[]M
  type W view[N]view[]M
  V.group(func(M)N) -> W
  (Note: the key type of the inner view remains unchanged)
*/
func Group(v IndexedView, f func(interface{}) interface{}) Deferred {
	return &group_op{src: v, fun: f}
}

/*
  type V view[]M
  type W view[N]M
  V.index(func(M)N) -> W
  (Note: the member values remain unchanged)

  v.index(f) is equivalent to v.group(f).collect(func(v view[]M)M{ return v.detect(func(_){return true}) })
*/
func Index(v IndexedView, f func(interface{}) interface{}) Deferred {
	return &index_op{src: v, fun: f}
}

/*
  type V view[]M
  V.sort(func(M)N) -> V
  (Note: the key type is lost)
*/
func Sort(v IndexedView, f func(interface{}) interface{}) Deferred {
	return &sort_op{src: v, fun: f}
}

func Union(v ...Deferred) Deferred {
	panic("not yet implemented")
}

type (
	table_op struct {
		name string
	}

	select_func func(interface{}) bool
	select_op   struct {
		src IndexedView
		fun select_func
	}

	reject_func func(interface{}) bool
	reject_op   struct {
		src IndexedView
		fun reject_func
	}

	collect_func func(interface{}) interface{}
	collect_op   struct {
		src IndexedView
		fun collect_func
	}

	group_func func(interface{}) interface{}
	group_op   struct {
		src IndexedView
		fun group_func
	}

	index_func func(interface{}) interface{}
	index_op   struct {
		src IndexedView
		fun index_func
	}

	sort_func func(interface{}) interface{}
	sort_op   struct {
		src IndexedView
		fun sort_func
	}
)

func (op *table_op) Resolve(txn *Transaction, events chan<- Event) {
	table := txn.GetTable(op.name)

	for _, change := range txn.changes {
		if change.Table != op.name {
			continue
		}

		switch change.Kind {
		case SET:
			old_kv_sha, new_kv_sha, changed := table.Set(change.Key, change.Value)
			if changed {
				events <- &ev_SET{op.name, old_kv_sha, new_kv_sha}
			}

		case UNSET:
			old_kv_sha, deleted := table.Del(change.Key)
			if deleted {
				events <- &ev_DEL{op.name, old_kv_sha}
			}

		}
	}

	old_tab_sha, new_tab_sha, changed := table.Commit()
	if changed {
		events <- &ev_CONSISTENT{op.name, old_tab_sha, new_tab_sha}
	}
}

func (op *select_op) Resolve(txn *Transaction, events chan<- Event)  {}
func (op *reject_op) Resolve(txn *Transaction, events chan<- Event)  {}
func (op *collect_op) Resolve(txn *Transaction, events chan<- Event) {}
func (op *group_op) Resolve(txn *Transaction, events chan<- Event)   {}
func (op *index_op) Resolve(txn *Transaction, events chan<- Event)   {}
func (op *sort_op) Resolve(txn *Transaction, events chan<- Event)    {}

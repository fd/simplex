package runtime

import (
	"simplex.sh/cas"
)

func DeclareTable(name string) Resolver {
	return &table_op{name}
}

/*
  type V view[]M
  V.select(func(M)bool) -> V
*/
func Select(v IndexedView, f select_func, name string) Resolver {
	return &select_op{src: v, fun: f, name: name}
}

/*
  type V view[]M
  V.reject(func(M)bool) -> V
*/
func Reject(v IndexedView, f reject_func, name string) Resolver {
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
func Collect(v IndexedView, f collect_func, name string) Resolver {
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
func Group(v IndexedView, f group_func, name string) Resolver {
	return &group_op{src: v, fun: f, name: name}
}

/*
  type V view[]M
  type W view[N]M
  V.index(func(M)N) -> W
  (Note: the member values remain unchanged)

  v.index(f) is equivalent to v.group(f).collect(func(v view[]M)M{ return v.detect(func(_){return true}) })
*/
func Index(v IndexedView, f index_func, name string) Resolver {
	return &index_op{src: v, fun: f, name: name}
}

/*
  type V view[]M
  V.sort(func(M)N) -> V
  (Note: the key type is lost)
*/
func Sort(v IndexedView, f sort_func, name string) Resolver {
	return &sort_op{src: v, fun: f, name: name}
}

func Union(v ...Resolver) Resolver {
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

	reject_func func(*Context, cas.Addr) bool
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

	group_func func(*Context, cas.Addr) interface{}
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

	sort_func func(*Context, cas.Addr) interface{}
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

func (op *index_op) Resolve(txn *Transaction) IChange { return IChange{} }

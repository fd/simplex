package data

import (
	"fmt"
	"github.com/fd/w/util"
)

var data_view_counters = map[string]int{}

type View struct {
	engine  *Engine
	current transformation
}

type transformation_info struct {
	Id          string
	Downstreams []string
}

func (i *transformation_info) Info() *transformation_info {
	return i
}

type transformation interface {
	Info() *transformation_info

	Transform(txn transaction)
}

type transformation_state interface {
	StoreReader

	Created() []string
	Updated() []string
	Destroyed() []string
}

func (v View) push(t transformation) View {
	ti := t.Info()
	ti.Id = v.new_id()

	if v.current != nil {
		ci := v.current.Info()
		ci.Downstreams = append(ci.Downstreams, ti.Id)
	}

	v.engine.transformations[ti.Id] = t

	v.current = t
	return v
}

func (v View) new_id() string {
	pkg := util.InitializingPackage()
	data_view_counters[pkg] += 1
	return fmt.Sprintf("%s:%d", pkg, data_view_counters[pkg])
}

/*
type GroupFunc func(Document) Value

func (v View) Group(f GroupFunc) View {
  v.group = f
  return v
}

func (v View) Paginate(n int) View {
  v.page = n
  return v
}
*/

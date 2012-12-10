package data

import (
	"fmt"
	"github.com/fd/w/util"
)

var data_view_counters = map[string]int{}

type View struct {
	controller *transformation_controller
}

type transformation_controller struct {
	id string

	prev           []*transformation_controller
	transformation transformation
	next           []*transformation_controller
}

type transformation interface {
	Transform(prev State, txn transaction)
}

func (v View) add_transformation(t Transformation) View {
	pkg := util.InitializingPackage()
	data_view_counters[pkg] += 1

	c := &transformation_controller{transformation: t}
	c.id = fmt.Sprintf("%s:%d", pkg, data_view_counters[pkg])

	c.prev = append(c.prev, v.controller)
	v.controller.next = append(v.controller.next, c)
	v.controller = c

	return v
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

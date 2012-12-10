package data

import (
	"fmt"
	"github.com/fd/w/util"
)

var data_view_counters = map[string]int{}

type View struct {
	engine  *Engine
	current *transformation_decl
}

type transformation_decl struct {
	id string

	upstream       []string
	transformation transformation
	downstream     []string
}

type transformation interface {
	Transform(txn transaction)
}

func (v View) add_transformation(t transformation) View {
	pkg := util.InitializingPackage()
	data_view_counters[pkg] += 1

	c := &transformation_decl{transformation: t}
	c.id = fmt.Sprintf("%s:%d", pkg, data_view_counters[pkg])

	fmt.Printf("View[%s]\n", c.id)

	// bind dependencies
	if v.current != nil {
		c.upstream = append(c.upstream, v.current.id)
		v.current.downstream = append(v.current.downstream, c.id)
	}

	// register with engine
	if v.engine.transformations == nil {
		v.engine.transformations = map[string]*transformation_decl{}
	}
	v.engine.transformations[c.id] = c

	// update view
	v.current = c
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

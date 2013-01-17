package data

import (
	"fmt"
	"github.com/fd/simplex/w/util"
)

type View struct {
	engine  *Engine
	current transformation
}

func (v View) push(t transformation) View {
	if v.current != nil {
		v.current.PushDownstream(t)
	}

	v.engine.transformations[t.Id()] = t
	v.current = t
	return v
}

func (v View) new_id() string {
	pkg := util.InitializingPackage()
	v.engine.transformation_counters[pkg] += 1
	return fmt.Sprintf("%s:%d", pkg, v.engine.transformation_counters[pkg])
}

/*
type GroupFunc func(Document) Value

func (v View) Group(f GroupFunc) GroupView {
  v.group = f
  return v
}

func (v View) Paginate(n int) PageView {
  v.page = n
  return v
}
*/

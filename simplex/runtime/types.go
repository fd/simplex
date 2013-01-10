package runtime

type View struct {
	Type  string
	Chain []string
}

type ViewWrapper interface {
	View() View
}

type GroupFunc func(m interface{}) interface{}
type SortFunc func(m interface{}) interface{}
type CollectFunc func(m interface{}) interface{}
type SelectFunc func(m interface{}) bool

func Source(typ string) View {
	return View{Type: typ}
}

func (v View) Select(f SelectFunc) View {
	return v
}
func (v View) Sort(f SortFunc) View {
	return v
}
func (v View) Group(f GroupFunc) View {
	return v
}
func (v View) Collect(f CollectFunc) View {
	return v
}

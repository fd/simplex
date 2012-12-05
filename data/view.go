package data

type SelectFunc func(Document) bool
type MapFunc func(Document) Value
type SortFunc func(Document) Value
type GroupFunc func(Document) Value

type View struct {
	selects []SelectFunc // Select on the raw input
	maps    []MapFunc    // Map the inputs
	sort    SortFunc
	group   GroupFunc

	offset int
	limit  int
	page   int

	state view_state
}

func (v View) Select(f SelectFunc) View {
	v.selects = append(v.selects, f)
	return v
}

func (v View) Map(f MapFunc) View {
	v.maps = append(v.maps, f)
	return v
}

func (v View) Sort(f SortFunc) View {
	v.sort = f
	return v
}

func (v View) Group(f GroupFunc) View {
	v.group = f
	return v
}

func (v View) Offset(n int) View {
	v.offset = n
	return v
}

func (v View) Limit(n int) View {
	v.limit = n
	return v
}

func (v View) Paginate(n int) View {
	v.page = n
	return v
}

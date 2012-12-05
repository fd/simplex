package data

type view_state struct {
	NextProductId int
	Products      map[int]Value
	Sorted        []sort_entry
	Groups        map[string][]int // map group key to set of product ids (sorted)
}

type sort_entry struct {
	Key Value // sort value
	Id  int   // product id
}

func (s *view_state) Add(v *View, d Document) {

	if !s.is_selected(v, d) {
		// not part of this collection
		return
	}

	// store the product
	value := s.map_document(v, d)
	s.NextProductId += 1
	id := s.NextProductId
	s.Products[id] = value

	// sort the products
	if v.sort != nil {
		s.Sorted = append(s.Sorted, sort_entry{v.sort(value), id})
		sort.Sort(s.Sorted)
	}

	// group the product
	if v.group != nil {
		group_id := v.group(val)
		l := s.Groups[group_id]
		l = append(l, id)
		sort.Sort(l)
		s.Groups[group_id] = l
	}
}

// test if a document is selected
func (s *view_state) is_selected(v *View, d Document) bool {
	for _, f := range v.selects {
		if !f(d) {
			return false
		}
	}
	return true
}

// map the document to a product
func (s *view_state) map_document(v *View, d Document) Value {
	var p Value = d
	for _, f := range v.maps {
		p = f(p)
	}
	return p
}

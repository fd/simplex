package data

/*
Group
*/

type group struct {
	u Collection
	f GroupFunc
	d Collection
	s group_state
}

type group_state struct {
	Products map[int]Value
	Sorted   []int
	Groups   map[string][]int
}

func (s *group_state) Update(changed_doc_ids []int) {
	// remove all products for changed_doc_ids
	// get products by id and add them back
}

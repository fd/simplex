package static

// Slice makes a new slice from a collection by taking only the objects between [beg,end)
func (in *C) Slice(beg, end int) *C {
	return in.Transform(in.elem_type, func(i_elems []interface{}) ([]interface{}, error) {
		if end > len(i_elems) {
			return i_elems[beg:len(i_elems)], nil
		}

		return i_elems[beg:end], nil
	})
}

// Frame() is a helper around Slice(). It works like the SQL LIMIT expression.
func (in *C) Frame(offset, limit int) *C {
	return in.Slice(offset, offset+limit)
}

package static

/*
Group puts the values of a collection in the groups named by the results of f.

  posts_by_tags := c.Group(func(v interface{})[]string{
    post := v.(*Post)
    return post.Tags
  })
*/
func (in *C) Group(f func(v interface{}) []string) *G {
	var (
		out = &G{elem_type: in.elem_type, tx: in.tx}
	)

	out.t.Do(func() error {
		if err := in.t.Wait(); err != nil {
			return err
		}

		var (
			groups = make(map[string][]interface{}, len(in.elems)/4)
		)

		for _, elem := range in.elems {
			for _, group_name := range f(elem) {
				groups[group_name] = append(groups[group_name], elem)
			}
		}

		out.elems = groups
		return nil
	})

	return out
}

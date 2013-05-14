package static

func (in *G) Get(group string) *C {
	var (
		out = &C{elem_type: in.elem_type, tx: in.tx}
	)

	out.t.Do(func() error {
		if err := in.t.Wait(); err != nil {
			return err
		}

		out.elems = in.elems[group]
		return nil
	})

	return out
}

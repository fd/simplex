package static

func (a *I) Add(b *I) *I {
	var (
		out = &I{elem_type: a.elem_type, tx: a.tx}
	)

	out.t.Do(func() error {
		if err := a.t.Wait(); err != nil {
			return err
		}

		if err := b.t.Wait(); err != nil {
			return err
		}

		var (
			mapped = make(map[string]interface{}, len(a.elems)+len(b.elems))
		)

		for k, v := range b.elems {
			mapped[k] = v
		}

		for k, v := range a.elems {
			mapped[k] = v
		}

		out.elems = mapped
		return nil
	})

	return out
}

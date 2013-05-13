package static

import (
	"simplex.sh/errors"
)

func (in *C) Index(f func(v interface{}) string) *I {
	var (
		out = &I{elem_type: in.elem_type, tx: in.tx}
	)

	out.t.Do(func() error {
		if err := in.t.Wait(); err != nil {
			return err
		}

		var (
			mapped = make(map[string]interface{}, len(in.elems))
		)

		for _, elem := range in.elems {
			key := f(elem)
			if _, p := mapped[key]; p {
				return errors.Fmt("duplicate key: %s", key)
			}
			mapped[key] = elem
		}

		out.elems = mapped
		return nil
	})

	return out
}

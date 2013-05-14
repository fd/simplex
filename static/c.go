package static

import (
	"reflect"
	"simplex.sh/future"
)

type C struct {
	t         future.Deferred
	tx        *Tx
	elem_type reflect.Type
	elems     []interface{}
}

func (c *C) ElemType() reflect.Type {
	return c.elem_type
}

func (c *C) Tx() *Tx {
	return c.tx
}

func (c *C) Wait() error {
	return c.t.Wait()
}

func (in *C) Transform(typ reflect.Type, f func(elems []interface{}) ([]interface{}, error)) *C {
	var (
		out = &C{elem_type: typ, tx: in.tx}
	)

	out.t.Do(func() error {
		if err := in.t.Wait(); err != nil {
			return err
		}

		elems, err := f(in.elems)
		if err != nil {
			return err
		}

		out.elems = elems
		return nil
	})

	return out
}

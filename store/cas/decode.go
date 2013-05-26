package cas

import (
	"encoding/gob"
)

func Decode(g Getter, addr Addr, v interface{}) error {
	r, err := g.Get(addr)
	if err != nil {
		return err
	}

	err = gob.NewDecoder(r).Decode(v)
	if err != nil {
		return err
	}

	return nil
}

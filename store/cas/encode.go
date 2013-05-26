package cas

import (
	"encoding/gob"
	"simplex.sh/store/digest"
)

func Encode(s Setter, v interface{}) (Addr, error) {
	addr, err := digest.Digest(v)
	if err != nil {
		return nil, err
	}

	w := s.Set()

	err = gob.NewEncoder(w).Encode(v)
	if err != nil {
		w.Abort()
		return nil, err
	}

	err = w.Commit(addr)
	if err != nil {
		return nil, err
	}

	return addr, nil
}

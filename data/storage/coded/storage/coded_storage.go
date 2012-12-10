package storage

import (
	"github.com/fd/w/data/storage/raw/driver"
)

type S struct {
	Driver driver.I
	Coder  Coder
}

func (s *S) Ids() ([]string, error) {
	return s.Driver.Ids()
}

func (s *S) Get(id string) (interface{}, error) {
	dat, err := s.Driver.Get(id)
	if err != nil {
		return nil, err
	}

	val, err := s.Coder.Decode(dat)
	if err != nil {
		return nil, err
	}

	return val, nil
}

func (s *S) Commit(set map[string]interface{}, del []string) error {
	n := make(map[string][]byte, len(set))

	for id, val := range set {
		dat, err := s.Coder.Encode(val)
		if err != nil {
			return err
		}

		n[id] = dat
	}

	return s.Driver.Commit(n, del)
}

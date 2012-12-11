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
	if len(dat) == 0 {
		return nil, nil
	}

	var val interface{}
	err = s.Coder.Decode(dat, &val)
	if err != nil {
		return nil, err
	}

	return val, nil
}

func (s *S) Restore(id string, val interface{}) error {
	dat, err := s.Driver.Get(id)
	if err != nil {
		return err
	}
	if len(dat) == 0 {
		return nil
	}

	err = s.Coder.Decode(dat, val)
	if err != nil {
		return err
	}

	return nil
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

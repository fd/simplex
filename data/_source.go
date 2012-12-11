package data

import (
	"github.com/fd/w/data/storage/coded/driver"
	"github.com/fd/w/data/storage/coded/prefixed"
	"github.com/fd/w/data/storage/coded/storage"
	raw "github.com/fd/w/data/storage/raw/driver"
)

type Source struct {
	driver driver.I
}

func NewSource(s raw.I) *Source {
	return &Source{
		driver: &prefixed.S{
			Prefix: "source/",
			Driver: &storage.S{
				Coder:  &storage.JsonCoder{},
				Driver: s,
			},
		},
	}
}

func (s *Source) Ids() []string {
	ids, err := s.driver.Ids()
	if err != nil {
		panic(err)
	}
	return ids
}

func (s *Source) Get(id string) Value {
	val, err := s.driver.Get(id)
	if err != nil {
		panic(err)
	}
	return Value(val)
}

func (s *Source) Commit(set map[string]Value, del []string) {
	n := make(map[string]interface{}, len(set))

	for id, val := range set {
		n[id] = interface{}(val)
	}

	err := s.driver.Commit(n, del)
	if err != nil {
		panic(err)
	}
}

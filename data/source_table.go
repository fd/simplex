package data

import (
	"github.com/fd/w/data/storage/coded/driver"
	"github.com/fd/w/data/storage/coded/prefixed"
	"github.com/fd/w/data/storage/coded/storage"
	raw "github.com/fd/w/data/storage/raw/driver"
)

type SourceTable struct {
	driver driver.I
}

func NewSourceTable(s raw.I) *SourceTable {
	return &SourceTable{
		driver: &prefixed.S{
			Prefix: "source/",
			Driver: &storage.S{
				Coder:  &storage.JsonCoder{},
				Driver: s,
			},
		},
	}
}

func (s *SourceTable) Ids() []string {
	ids, err := s.driver.Ids()
	if err != nil {
		panic(err)
	}
	return ids
}

func (s *SourceTable) Get(id string) Value {
	val, err := s.driver.Get(id)
	if err != nil {
		panic(err)
	}
	return Value(val)
}

func (s *SourceTable) Commit(set map[string]Value, del []string) {
	n := make(map[string]interface{}, len(set))

	for id, val := range set {
		n[id] = interface{}(val)
	}

	err := s.driver.Commit(n, del)
	if err != nil {
		panic(err)
	}
}

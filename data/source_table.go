package data

import (
	"github.com/fd/w/data/storage/coded/driver"
	"github.com/fd/w/data/storage/coded/prefixed"
	"github.com/fd/w/data/storage/coded/storage"
	raw "github.com/fd/w/data/storage/raw/driver"
)

type source_table struct {
	driver driver.I
	ids    []string
}

func new_source_table(s raw.I) *source_table {
	return &source_table{
		driver: &prefixed.S{
			Prefix: "source/",
			Driver: &storage.S{
				Coder:  &storage.JsonCoder{},
				Driver: s,
			},
		},
	}
}

func (s *source_table) Ids() []string {
	if len(s.ids) != 0 {
		return s.ids
	}

	ids, err := s.driver.Ids()
	if err != nil {
		panic(err)
	}

	s.ids = ids
	return ids
}

func (s *source_table) Get(id string) Value {
	val, err := s.driver.Get(id)
	if err != nil {
		panic(err)
	}
	return Value(val)
}

func (s *source_table) Commit(set map[string]Value, del []string) {
	n := make(map[string]interface{}, len(set))

	for id, val := range set {
		n[id] = interface{}(val)
	}

	err := s.driver.Commit(n, del)
	if err != nil {
		panic(err)
	}
}

package data

import (
	"github.com/fd/w/data/storage/coded/driver"
	"github.com/fd/w/data/storage/coded/prefixed"
	"github.com/fd/w/data/storage/coded/storage"
	raw "github.com/fd/w/data/storage/raw/driver"
)

type StateTable struct {
	driver driver.I
}

func NewStateTable(s raw.I) *StateTable {
	return &StateTable{
		driver: &prefixed.S{
			Prefix: "state/",
			Driver: &storage.S{
				Coder:  &storage.GobCoder{},
				Driver: s,
			},
		},
	}
}

func (s *StateTable) Ids() []string {
	ids, err := s.driver.Ids()
	if err != nil {
		panic(err)
	}
	return ids
}

func (s *StateTable) Restore(id string, state interface{}) {
	err := s.driver.Restore(id, state)
	if err != nil {
		panic(err)
	}
}

func (s *StateTable) Commit(set map[string]Value, del []string) {
	n := make(map[string]interface{}, len(set))

	for id, val := range set {
		n[id] = interface{}(val)
	}

	err := s.driver.Commit(n, del)
	if err != nil {
		panic(err)
	}
}

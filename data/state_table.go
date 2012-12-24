package data

import (
	"github.com/fd/w/data/storage"
)

type state_table struct {
	driver driver.I
}

func new_state_table(s raw.I) *state_table {
	return &state_table{
		driver: &prefixed.S{
			Prefix: "state/",
			Driver: &storage.S{
				Coder:  &storage.GobCoder{},
				Driver: s,
			},
		},
	}
}

func (s *state_table) Ids() []string {
	ids, err := s.driver.Ids()
	if err != nil {
		panic(err)
	}
	return ids
}

func (s *state_table) Restore(id string, state interface{}) {
	err := s.driver.Restore(id, state)
	if err != nil {
		panic(err)
	}
}

func (s *state_table) Commit(set map[string]Value, del []string) {
	n := make(map[string]interface{}, len(set))

	for id, val := range set {
		n[id] = interface{}(val)
	}

	err := s.driver.Commit(n, del)
	if err != nil {
		panic(err)
	}
}

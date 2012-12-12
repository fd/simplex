package data

import (
	"github.com/fd/w/data/storage/coded/driver"
	"github.com/fd/w/data/storage/coded/prefixed"
	"github.com/fd/w/data/storage/coded/storage"
	raw "github.com/fd/w/data/storage/raw/driver"
	"net/http"
)

type target_table struct {
	driver driver.I
}

func new_target_table(s raw.I) *target_table {
	return &target_table{
		driver: &prefixed.S{
			Prefix: "target/",
			Driver: &storage.S{
				Coder:  &storage.GobCoder{},
				Driver: s,
			},
		},
	}
}

func (s *target_table) Ids() []string {
	ids, err := s.driver.Ids()
	if err != nil {
		panic(err)
	}
	return ids
}

func (s *target_table) Get(id string) Artefact {
	val, err := s.driver.Get(id)
	if err != nil {
		panic(err)
	}
	if art, ok := val.(Artefact); ok {
		return art
	}
	panic("unable to convert to Artefact")
}

func (s *target_table) Commit(set map[string]Artefact, del []string) {
	n := make(map[string]interface{}, len(set))

	for id, val := range set {
		n[id] = interface{}(val)
	}

	err := s.driver.Commit(n, del)
	if err != nil {
		panic(err)
	}
}

type Artefact struct {
	Header http.Header
	Body   []byte
}

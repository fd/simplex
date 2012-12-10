package data

import (
	"github.com/fd/w/data/storage/coded/driver"
	"github.com/fd/w/data/storage/coded/prefixed"
	"github.com/fd/w/data/storage/coded/storage"
	raw "github.com/fd/w/data/storage/raw/driver"
	"net/http"
)

type Target struct {
	driver driver.I
}

func NewTarget(s raw.I) *Target {
	return &Target{
		driver: &prefixed.S{
			Prefix: "target/",
			Driver: &storage.S{
				Coder:  &storage.GobCoder{},
				Driver: s,
			},
		},
	}
}

func (s *Target) Ids() []string {
	ids, err := s.driver.Ids()
	if err != nil {
		panic(err)
	}
	return ids
}

func (s *Target) Get(id string) Artefact {
	val, err := s.driver.Get(id)
	if err != nil {
		panic(err)
	}
	if art, ok := val.(Artefact); ok {
		return art
	}
	panic("unable to convert to Artefact")
}

func (s *Target) Commit(set map[string]Artefact, del []string) {
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

package prefixed

import (
	"github.com/fd/w/data/storage/coded/driver"
	"strings"
)

type S struct {
	Driver driver.I
	Prefix string
}

func (s *S) Ids() ([]string, error) {
	all, err := s.Driver.Ids()
	if err != nil {
		return nil, err
	}

	ids := make([]string, 0, 1024)
	l := len(s.Prefix)

	for _, id := range all {
		if strings.HasPrefix(id, s.Prefix) {
			ids = append(ids, id[0:l])
		}
	}

	return ids, nil
}

func (s *S) Get(id string) (interface{}, error) {
	return s.Driver.Get(s.Prefix + id)
}

func (self *S) Commit(set map[string]interface{}, del []string) error {
	s := make(map[string]interface{}, len(set))
	d := make([]string, len(del))

	for id, val := range set {
		s[self.Prefix+id] = val
	}

	for i, id := range del {
		d[i] = self.Prefix + id
	}

	return self.Driver.Commit(s, d)
}

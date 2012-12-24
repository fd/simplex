package memory

import (
	"github.com/fd/w/data/ident"
	"github.com/fd/w/data/storage/driver"
	"net/url"
)

func init() {
	driver.Register("mem", func(us string) (driver.I, error) {
		_, err := url.Parse(us)
		if err != nil {
			return nil, err
		}

		return &S{objects: map[string][]byte{}}, nil
	})
}

type S struct {
	objects map[string][]byte
}

func (s *S) Get(key ident.SHA) ([]byte, error) {
	hex := ident.HexSHA(key)
	dat, found := s.objects[hex]
	if !found {
		return nil, driver.NotFound
	}
	return dat, nil
}

func (s *S) Set(key ident.SHA, val []byte) error {
	hex := ident.HexSHA(key)
	s.objects[hex] = val
	return nil
}

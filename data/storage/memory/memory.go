package memory

import (
	"encoding/hex"
	"github.com/fd/simplex/data/storage/driver"
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

func (s *S) Get(key [20]byte) ([]byte, error) {
	hex := hex.EncodeToString(key[:])
	dat, found := s.objects[hex]
	if !found {
		return nil, driver.NotFound
	}
	return dat, nil
}

func (s *S) Set(key [20]byte, val []byte) error {
	hex := hex.EncodeToString(key[:])
	s.objects[hex] = val
	return nil
}

package redis

import (
	"encoding/hex"
	"fmt"
	"github.com/fd/simplex/data/storage/driver"
	"github.com/simonz05/godis/redis"
	"net/url"
)

func init() {
	driver.Register("redis", func(us string) (driver.I, error) {
		u, err := url.Parse(us)
		if err != nil {
			return nil, err
		}

		fmt.Printf("REDIS: %s\n", u.Host)
		c := redis.New("tcp:"+u.Host, 0, "")
		return &S{Client: c}, nil
	})
}

type S struct {
	Client *redis.Client
}

func (s *S) Get(key [20]byte) ([]byte, error) {
	key_str := hex.EncodeToString([]byte(key[:]))

	elem, err := s.Client.Get(key_str)
	if err != nil {
		return nil, err
	}

	if len(elem) == 0 {
		return nil, driver.NotFound
	}

	return elem, nil
}

func (s *S) Set(key [20]byte, val []byte) error {
	key_str := hex.EncodeToString([]byte(key[:]))

	return s.Client.Set(key_str, val)
}

package file_system

import (
	"encoding/hex"
	"github.com/fd/simplex/data/storage/driver"
	"io/ioutil"
	"net/url"
	"os"
	"path"
)

func init() {
	driver.Register("file", func(us string) (driver.I, error) {
		u, err := url.Parse(us)
		if err != nil {
			return nil, err
		}

		return &S{Root: u.Path}, nil
	})
}

type S struct {
	Root string
}

func (s *S) Get(key [20]byte) ([]byte, error) {
	var (
		pat = s.path_for_sha(key)
	)

	dat, err := ioutil.ReadFile(pat)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, driver.NotFound
		}
		return nil, err
	}
	return dat, nil
}

func (s *S) Set(key [20]byte, val []byte) error {
	var (
		pat = s.path_for_sha(key)
	)

	_, err := os.Stat(pat)
	if err == nil {
		return nil
	} else {
		if os.IsNotExist(err) {
			err = nil
		} else {
			return err
		}
	}

	err = os.MkdirAll(path.Dir(pat), 0755)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(pat, val, 0644)
}

func (s *S) path_for_sha(sha [20]byte) string {
	hex := hex.EncodeToString(sha[:])
	a := hex[0:4]
	b := hex[4:4]
	c := hex[8:]
	return path.Join(s.Root, "objects", a, b, c)
}

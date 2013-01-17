package file_system

import (
	"github.com/fd/w/data/ident"
	"github.com/fd/w/data/storage/driver"
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

func (s *S) Get(key ident.SHA) ([]byte, error) {
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

func (s *S) Set(key ident.SHA, val []byte) error {
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

func (s *S) path_for_sha(sha ident.SHA) string {
	hex := ident.HexSHA(sha)
	a := hex[0:4]
	b := hex[4:4]
	c := hex[8:]
	return path.Join(s.Root, "objects", a, b, c)
}

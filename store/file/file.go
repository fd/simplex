package file

import (
	"github.com/fd/static/errors"
	"github.com/fd/static/store"
	"io"
	"net/url"
	"os"
	"path"
)

func init() {
	store.Register("file", Open)
}

func Open(u *url.URL) (store.Store, error) {
	if u.Host != "" {
		wd, err := os.Getwd()
		if err != nil {
			return nil, errors.Forward(err, "Failed to open store at `%s`", u)
		}

		u.Path = path.Join(wd, u.Host+u.Path)
		u.Host = ""
	}
	return file_store_t(u.Path), nil
}

type file_store_t string

func (f file_store_t) Get(name string) (io.ReadCloser, error) {
	r, err := os.Open(path.Join(string(f), name))

	if err != nil {
		if os.IsNotExist(err) {
			err = store.NotFoundError(name)
		} else {
			err = errors.Forward(err, "Unable to Get() `%s`", name)
		}
	}

	return r, err
}

func (f file_store_t) Set(name string) (io.WriteCloser, error) {
	p := path.Join(string(f), name)

	err := os.MkdirAll(path.Dir(p), 0755)
	if err != nil {
		return nil, errors.Forward(err, "Unable to Set() `%s`", name)
	}

	w, err := os.Create(p)

	if err != nil {
		err = errors.Forward(err, "Unable to Set() `%s`", name)
	}

	return w, err
}

package file_system

import (
	"github.com/fd/w/data/storage/raw"
	"github.com/fd/w/data/storage/raw/driver"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"strings"
)

func init() {
	raw.Register("file", func(us string) (driver.I, error) {
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

func (f *S) Ids() ([]string, error) {
	err := os.MkdirAll(f.Root, 0755)
	if err != nil {
		return nil, err
	}

	return f.ids_for_dir(f.Root, "")
}

func (f *S) ids_for_dir(dir, prefix string) ([]string, error) {
	d, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer d.Close()

	fis, err := d.Readdir(-1)
	if err != nil {
		return nil, err
	}

	ids := make([]string, 0, len(fis))
	for _, fi := range fis {
		name := fi.Name()

		if name[:1] == "." || name[:1] == "_" {
			continue
		}

		if fi.IsDir() {
			pref := name
			if prefix != "" {
				pref = path.Join(prefix, name)
			}
			i, err := f.ids_for_dir(path.Join(dir, name), pref)
			if err != nil {
				return nil, err
			}
			ids = append(ids, i...)
			continue
		}

		if !strings.HasSuffix(name, ".dat") {
			continue
		}

		id := name[:len(name)-4]
		ids = append(ids, path.Join(prefix, id))
	}

	return ids, nil
}

func (f *S) Get(id string) ([]byte, error) {
	pat := path.Join(f.Root, id)
	dir := path.Dir(pat)

	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(pat + ".dat")
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return data, nil
}

func (f *S) Commit(set map[string][]byte, del []string) error {
	err := os.MkdirAll(f.Root, 0755)
	if err != nil {
		return err
	}

	for _, id := range del {
		err := os.Remove(path.Join(f.Root, id+".dat"))
		if err != nil {
			if os.IsNotExist(err) {
				err = nil
			} else {
				return err
			}
		}
	}

	for id, data := range set {
		pat := path.Join(f.Root, id+".dat")
		os.MkdirAll(path.Dir(pat), 0755)
		err := ioutil.WriteFile(pat, data, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

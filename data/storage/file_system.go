package storage

import (
	"encoding/hex"
	"io/ioutil"
	"os"
	"path"
)

type FileSystem struct {
	Root string
}

func (f *FileSystem) Ids() ([]string, error) {
	err := os.MkdirAll(f.Root, 0755)
	if err != nil {
		return nil, err
	}

	d, err := os.Open(f.Root)
	if err != nil {
		return nil, err
	}
	defer d.Close()

	names, err := d.Readdirnames(-1)
	if err != nil {
		return nil, err
	}

	ids := make([]string, 0, len(name))
	for _, name := range names {
		if name[:1] == "." || name[:1] == "_" {
			continue
		}

		if !strings.HasSuffix(name, ".dat") {
			continue
		}

		id := name[:len(name)-4]
		id, err := hex.DecodeString(id)
		if err != nil {
			return nil, err
		}

		ids = append(ids, id)
	}

	return ids, nil
}

func (f *FileSystem) Get(id string) ([]byte, error) {
	id = hex.EncodeToString([]byte(id))

	err := os.MkdirAll(f.Root, 0755)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(path.Join(f.Root, id+".dat"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return data, nil
}

func (f *FileSystem) Commit(set map[string][]byte, del []string) error {
	err := os.MkdirAll(f.Root, 0755)
	if err != nil {
		return err
	}

	for _, id := range del {
		id = hex.EncodeToString([]byte(id))

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
		id = hex.EncodeToString([]byte(id))

		err := ioutil.WriteFile(path.Join(f.Root, id+".dat"), data, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

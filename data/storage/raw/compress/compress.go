package compress

import (
	"bytes"
	"compress/zlib"
	"github.com/fd/w/data/storage/raw/driver"
	"io/ioutil"
)

type S struct {
	Storage driver.I
}

func (f *S) Ids() ([]string, error) {
	return f.Storage.Ids()
}

func (f *S) Get(id string) ([]byte, error) {
	dat, err := f.Storage.Get(id)
	if err != nil {
		return nil, err
	}
	if len(dat) == 0 {
		return nil, nil
	}

	var (
		buf_i *bytes.Buffer
	)

	buf_i = bytes.NewBuffer(dat)

	r, err := zlib.NewReader(buf_i)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	return ioutil.ReadAll(r)
}

func (f *S) Commit(set map[string][]byte, del []string) error {
	n := make(map[string][]byte, len(set))

	if set != nil {
		for k, v := range set {
			c, err := f.compress(v)
			if err != nil {
				return err
			}
			n[k] = c
		}
	}

	return f.Storage.Commit(n, del)
}

func (f *S) compress(dat []byte) ([]byte, error) {
	if len(dat) == 0 {
		return nil, nil
	}

	var (
		buf_o bytes.Buffer
	)

	w, err := zlib.NewWriterLevel(&buf_o, zlib.DefaultCompression)
	if err != nil {
		return nil, err
	}

	w.Write(dat)
	w.Close()

	return buf_o.Bytes(), nil
}

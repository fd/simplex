package storage

import (
	"bytes"
	"compress/zlib"
	"encoding/gob"
	"github.com/fd/w/data/ident"
	"github.com/fd/w/data/storage/driver"
)

type S struct {
	d driver.I
}

func New(url string) (*S, error) {
	d, err := driver.NewDriver(url)
	if err != nil {
		return nil, err
	}
	return &S{d}, nil
}

func (s *S) Get(key ident.SHA, val interface{}) (found bool) {
	data, err := s.d.Get(key)
	if err == driver.NotFound {
		return false
	}
	if err != nil {
		panic(err)
	}

	comp, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	defer comp.Close()

	err = gob.NewDecoder(comp).Decode(val)
	if err != nil {
		panic(err)
	}

	return true
}

func (s *S) Set(val interface{}) ident.SHA {
	var (
		sha = ident.Hash(val)
		buf bytes.Buffer
	)

	comp, err := zlib.NewWriterLevel(&buf, zlib.DefaultCompression)
	if err != nil {
		panic(err)
	}

	err = gob.NewEncoder(comp).Encode(val)
	if err != nil {
		panic(err)
	}

	comp.Close()

	err = s.d.Set(sha, buf.Bytes())
	if err != nil {
		panic(err)
	}

	return sha
}

func (s *S) SetEntry(key ident.SHA) {
	err := s.d.Set(ident.ZeroSHA, []byte(key[:]))
	if err != nil {
		panic(err)
	}
}

func (s *S) GetEntry() (sha ident.SHA, found bool) {
	b, err := s.d.Get(ident.ZeroSHA)
	if err == driver.NotFound {
		return ident.ZeroSHA, false
	}
	if err != nil {
		panic(err)
	}

	sha = ident.SHA{}
	copy(sha[:], b)
	return sha, true
}

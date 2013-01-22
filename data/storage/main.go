package storage

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"github.com/fd/simplex/data/blob"
	"github.com/fd/simplex/data/storage/driver"
	"io"
)

type SHA [20]byte

var ZeroSHA = SHA{}

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

func (s *S) Get(key SHA, val interface{}) (found bool) {
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

	err = blob.NewDecoder(comp).Decode(val)
	if err != nil {
		panic(err)
	}

	return true
}

func (s *S) Set(val interface{}) SHA {
	var (
		sha_w = sha1.New()
		sha   = SHA{}
		buf   bytes.Buffer
	)

	comp, err := zlib.NewWriterLevel(&buf, zlib.DefaultCompression)
	if err != nil {
		panic(err)
	}

	err = blob.NewEncoder(io.MultiWriter(comp, sha_w)).Encode(val)
	if err != nil {
		panic(err)
	}

	comp.Close()
	copy(sha[:], sha_w.Sum(nil))

	err = s.d.Set(sha, buf.Bytes())
	if err != nil {
		panic(err)
	}

	return sha
}

func (s *S) SetEntry(key SHA) {
	err := s.d.Set(ZeroSHA, []byte(key[:]))
	if err != nil {
		panic(err)
	}
}

func (s *S) GetEntry() (sha SHA, found bool) {
	b, err := s.d.Get(ZeroSHA)
	if err == driver.NotFound {
		return ZeroSHA, false
	}
	if err != nil {
		panic(err)
	}

	sha = SHA{}
	copy(sha[:], b)
	return sha, true
}

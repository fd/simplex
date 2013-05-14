package store

import (
	"bytes"
	"io"
	"io/ioutil"
)

type cache_store struct {
	store   Store
	entries map[string][]byte
}

func Cache(sub Store) Store {
	return &cache_store{
		store:   sub,
		entries: map[string][]byte{},
	}
}

func (c *cache_store) Set(name string) (io.WriteCloser, error) {
	return c.store.Set(name)
}

func (c *cache_store) Get(name string) (io.ReadCloser, error) {
	if data, p := c.entries[name]; p {
		return ioutil.NopCloser(bytes.NewReader(data)), nil
	}

	var (
		buf      = &bytes.Buffer{}
		upstream io.ReadCloser
		err      error
	)

	upstream, err = c.store.Get(name)
	if err != nil {
		return upstream, err
	}

	return &flush_closer{
		Reader: io.TeeReader(upstream, buf),
		b:      buf,
		u:      upstream,
		c:      c,
		n:      name,
	}, nil
}

type flush_closer struct {
	io.Reader
	b *bytes.Buffer
	u io.ReadCloser
	c *cache_store
	n string
}

func (f *flush_closer) Close() error {
	err := f.u.Close()

	if err == nil {
		f.c.entries[f.n] = f.b.Bytes()
	}

	return err
}

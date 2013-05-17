package redis

import (
	"bytes"
	"github.com/simonz05/godis/redis"
	"io"
	"io/ioutil"
	"net/url"
	"path"
	"simplex.sh/store"
)

func init() {
	store.Register("redis", Open)
}

type store_t struct {
	conn   *redis.Client
	prefix string
}

type setter struct {
	io.Writer
	conn *redis.Client
	key  string
}

func Open(u *url.URL) (store.Store, error) {
	var (
		addr   = "tcp:" + u.Host
		pass   = ""
		prefix = "sx:"
		conn   *redis.Client
	)

	if ui := u.User; ui != nil {
		pass, _ = ui.Password()
	}

	prefix += u.Path
	prefix = path.Clean(prefix)

	conn = redis.New(addr, 0, pass)

	return &store_t{conn, prefix}, nil
}

func (s *store_t) GetBlob(name string) (io.ReadCloser, error) {
	key := path.Join(s.prefix, name)
	elem, err := s.conn.Get(key)

	if err != nil {
		return nil, err
	}

	if elem == nil {
		return nil, store.NotFoundError(name)
	}

	return ioutil.NopCloser(bytes.NewReader(elem.Bytes())), err
}

func (s *store_t) SetBlob(name string) (io.WriteCloser, error) {
	key := path.Join(s.prefix, name)
	return &setter{bytes.NewBuffer(nil), s.conn, key}, nil
}

func (s *setter) Close() error {
	return s.conn.Set(s.key, s.Writer.(*bytes.Buffer).Bytes())
}

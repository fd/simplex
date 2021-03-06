package store

import (
	"net/url"
	"simplex.sh/errors"
)

type (
	Factory func(*url.URL) (Store, error)
)

var (
	registry = map[string]Factory{}
)

func Register(name string, f Factory) {
	if _, p := registry[name]; p {
		panic("duplicate store factory named: " + name)
	}

	registry[name] = f
}

func OpenOld(source string) (Store, error) {
	u, err := url.Parse(source)
	if err != nil {
		return nil, err
	}

	factory, p := registry[u.Scheme]
	if !p {
		return nil, errors.Fmt("Unknown store type: %s", u.Scheme)
	}

	return factory(u)
}

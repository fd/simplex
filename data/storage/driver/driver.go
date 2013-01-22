package driver

import (
	"errors"
	"fmt"
	"strings"
)

type I interface {
	Get(key [20]byte) ([]byte, error)
	Set(key [20]byte, val []byte) error
}

var NotFound = errors.New("Object not found")

type FactoryFunc func(url string) (I, error)

var factories = map[string]FactoryFunc{}

func Register(name string, factory FactoryFunc) {
	factories[name] = factory
}

func NewDriver(url string) (I, error) {
	parts := strings.SplitN(url, "://", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid storage url: %s", url)
	}

	factory, present := factories[parts[0]]
	if !present {
		return nil, fmt.Errorf("invalid storage type: %s", url)
	}

	return factory(url)
}

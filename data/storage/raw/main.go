package raw

import (
	"fmt"
	"github.com/fd/w/data/storage/raw/driver"
	"strings"
)

type Factory func(url string) (driver.I, error)

var factories = map[string]Factory{}

func Register(name string, factory Factory) {
	factories[name] = factory
}

func New(url string) (driver.I, error) {
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

package container

import (
	"os"
	"simplex.sh/errors"
	"simplex.sh/shttp"
	"simplex.sh/store"
	"sync"
)

var (
	app_registry []Factory
)

type Environment struct {
	Source      string
	Destination string
	HttpAddr    string
}

type container_t struct {
	env Environment

	src     store.Store
	dst     store.Store
	apps    []*Application
	app_map map[string]*Application
	router  shttp.HostRouter

	shutdown chan os.Signal
	mtx      sync.RWMutex
	err      errors.List
}

func new_container(env Environment) (*container_t, error) {

	if env.HttpAddr == "" {
		env.HttpAddr = ":3000"
	}

	src, err := store.Open(env.Source)
	if err != nil {
		return nil, err
	}

	dst, err := store.Open(env.Destination)
	if err != nil {
		return nil, err
	}

	c := &container_t{
		env:      env,
		app_map:  map[string]*Application{},
		src:      src,
		dst:      dst,
		shutdown: make(chan os.Signal, 1),
	}

	for _, factory := range app_registry {
		factory.build(c)
	}

	// reset the host router
	c.router.Reset()

	return c, nil
}

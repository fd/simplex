package container

import (
	"database/sql"
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
	Database string
	Source   string
	HttpAddr string
}

type container_t struct {
	env Environment

	database *sql.DB
	src      store.Store
	apps     []*Application
	app_map  map[string]*Application
	router   shttp.HostRouter

	shutdown chan os.Signal
	mtx      sync.RWMutex
	err      errors.List
}

func new_container(env Environment) (*container_t, error) {

	if env.HttpAddr == "" {
		env.HttpAddr = ":3000"
	}

	db, err := store.Open(env.Database)
	if err != nil {
		return nil, err
	}

	src, err := store.OpenOld(env.Source)
	if err != nil {
		return nil, err
	}

	c := &container_t{
		env:      env,
		app_map:  map[string]*Application{},
		database: db,
		src:      src,
		shutdown: make(chan os.Signal, 1),
	}

	for _, factory := range app_registry {
		err := factory.build(c)
		if err != nil {
			return nil, err
		}
	}

	// reset the host router
	c.router.Reset()

	return c, nil
}

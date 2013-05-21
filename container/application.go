package container

import (
	"database/sql"
	"github.com/gorilla/mux"
	"net/http"
	"simplex.sh/errors"
	"simplex.sh/shttp"
	"simplex.sh/static"
	"simplex.sh/store"
)

type Factory func(*Application)

func App(f Factory) Factory {
	app_registry = append(app_registry, f)
	return f
}

func (f Factory) build(c *container_t) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	app := &Application{
		container: c,
		dynamic:   mux.NewRouter(),
	}

	f(app)

	if app.Name == "" {
		return errors.Fmt("application name must not be empty.")
	}

	if _, p := c.app_map[app.Name]; p {
		return errors.Fmt("%s: application name must be unique.", app.Name)
	}

	if app.Generator == nil {
		return errors.Fmt("%s: application must have a generator.", app.Name)
	}

	app.database = c.database
	app.src = store.SubStore(c.src, app.Name)
	app.dst = store.SubStore(c.dst, app.Name)

	static, err := shttp.NewRouteHandler(store.Cache(app.dst))
	if err != nil {
		return errors.Forward(err, "%s: error while loading the route handler.", app.Name)
	}

	app.static = static
	app.dynamic.NotFoundHandler = static

	c.router.Add(app)
	c.apps = append(c.apps, app)
	c.app_map[app.Name] = app

	return nil
}

type Application struct {
	Name       string
	Generator  static.Generator
	ExtraHosts []string // extra hosts that route to this application

	container *container_t
	database  *sql.DB
	src       store.Store
	dst       store.Store
	dynamic   *mux.Router
	static    *shttp.RouteHandler
}

func (a *Application) Router() *mux.Router {
	return a.dynamic
}

func (a *Application) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.dynamic.ServeHTTP(w, r)
}

func (a *Application) Generate() error {
	return static.Generate(
		a.src,
		a.dst,
		a.database,
		a.Generator,
	)
}

func (a *Application) Hostnames() []string {
	var hosts []string

	for _, host := range a.ExtraHosts {
		hosts = append(hosts, host)
	}

	for _, host := range a.static.Hostnames() {
		hosts = append(hosts, host)
	}

	return hosts
}

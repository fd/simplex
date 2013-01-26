package api

import (
	"fmt"
	"github.com/fd/simplex/runtime"
	"net/http"
)

type (
	API struct {
		name   string
		env    *runtime.Environment
		tables map[string]runtime.Table
		views  map[string]runtime.IndexedView
		routes map[string]string
	}
)

func New(env *runtime.Environment, name string) *API {
	api := &API{
		name,
		env,
		map[string]runtime.Table{},
		map[string]runtime.IndexedView{},
		map[string]string{},
	}

	env.RegisterTerminal(api)

	return api
}

/*
  Register a table API endpoint.
  Available HTTP calls:

    GET    /_api                Return information on the current transaction.
    PATCH  /_api                Execute a list of actions in a single transaction.
    GET    /_api/{route}.json   Return a JSON document containign all the table entries
*/
func (api *API) RegisterTable(table runtime.Table) {
	if _, p := api.tables[table.TableId()]; p {
		panic(fmt.Sprintf("Already registered a table by the name of `%s`", table.TableId()))
	}

	api.tables[table.TableId()] = table
}

func (api *API) RegisterView(view runtime.IndexedView, route string) {
	if _, p := api.routes[route]; p {
		panic(fmt.Sprintf("Already registered a view at `%s`", route))
	}

	if _, p := api.views[view.DeferredId()]; p {
		panic(fmt.Sprintf("Already registered a view by the name of `%s`", view.DeferredId()))
	}

	api.routes[route] = view.DeferredId()
	api.views[view.DeferredId()] = view
}

func (api *API) DeferredId() string {
	return "API/" + api.name
}

func (api *API) Resolve(txn *runtime.Transaction, events chan<- runtime.Event) {
	var (
		funnel runtime.Funnel
	)

	for _, table := range api.tables {
		funnel.Add(txn.Resolve(table))
	}

	for _, view := range api.views {
		funnel.Add(txn.Resolve(view))
	}

	for _ = range funnel.Run() {
	}

}

func (api *API) ServeHTTP(w http.ResponseWriter, req *http.Request) {
}

/*
  Should return the current transactions SHA.
  This SHA MUST be passed allong with all the other requests in order to guarantee the consistency of the data.
*/
func (api *API) handle_GET_api(w http.ResponseWriter, req *http.Request) {

}

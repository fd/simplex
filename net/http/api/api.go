package api

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"simplex.sh/cas"
	"simplex.sh/cas/btree"
	"simplex.sh/runtime"
	"simplex.sh/runtime/event"
	"simplex.sh/runtime/promise"
	"strings"
)

type (
	API struct {
		name   string
		env    *runtime.Environment
		tables map[string]runtime.Table
		views  map[string]promise.Deferred
		routes map[string]string

		ViewTables map[string]*table_handle
	}

	table_handle struct {
		addr  cas.Addr
		table *btree.Tree
	}
)

func New(env *runtime.Environment, name string) *API {
	api := &API{
		name,
		env,
		map[string]runtime.Table{},
		map[string]promise.Deferred{},
		map[string]string{},
		map[string]*table_handle{},
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
	{ // normalize the route
		l := len(route)
		if route[l-1] == '/' {
			route = route[:l-1]
			l = len(route)
		}
		if l > 0 && route[0] != '/' {
			route = "/" + route
		}
	}

	if _, p := api.routes[route]; p {
		panic(fmt.Sprintf("Already registered a view at `%s`", route))
	}

	if _, p := api.views[view.DeferredId()]; p {
		panic(fmt.Sprintf("Already registered a view by the name of `%s`", view.DeferredId()))
	}

	sha := sha1.New()

	json_view := runtime.Collect(
		view,
		func(ctx *runtime.Context, m_addr cas.Addr) cas.Addr {
			var m reflect.Value
			m = reflect.New(view.EltType())
			ctx.LoadValue(m_addr, m)

			sha.Reset()
			sha.Write([]byte(m_addr[:]))

			data, err := json.Marshal(map[string]interface{}{
				"Vsn": hex.EncodeToString(sha.Sum(nil)),
				"Obj": m.Interface(),
			})
			if err != nil {
				panic(err)
			}

			return ctx.Save(&data)
		},
		"API/FORMAT_JSON/"+view.DeferredId(),
	)

	api.routes[route] = view.DeferredId()
	api.views[view.DeferredId()] = json_view
}

func (api *API) DeferredId() string {
	return "API/" + api.name
}

func (api *API) Resolve(state promise.State, events chan<- event.Event) {
	var (
		funnel event.Funnel
	)

	for _, table := range api.tables {
		funnel.Add(state.Resolve(table).C)
	}

	for _, view := range api.views {
		funnel.Add(state.Resolve(view).C)
	}

	for e := range funnel.Run() {
		// propagate error events
		if err, ok := e.(event.Error); ok {
			events <- err
			continue
		}

		event, ok := e.(*runtime.ConsistentTable)
		if !ok {
			continue
		}

		if strings.HasPrefix(event.Table, "API/FORMAT_JSON/") {
			name := event.Table[len("API/FORMAT_JSON/"):]

			if event.B == nil {
				delete(api.ViewTables, name)
			} else {
				api.ViewTables[name] = &table_handle{addr: event.B}
			}
		}
	}

}

func (api *API) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":

		if req.URL.Path == "/" {
			api.handle_GET_info(w, req)
			return
		}

		if view_id, p := api.routes[req.URL.Path]; p {
			if handle, p := api.ViewTables[view_id]; p {
				api.handle_GET_view(w, req, handle)
				return
			}
		}

	case "PATCH":

		if req.URL.Path == "/" {
			api.handle_PATCH_transaction(w, req)
			return
		}

	}

	http.NotFound(w, req)
}

/*
  Should return the current transactions SHA.
*/
func (api *API) handle_GET_info(w http.ResponseWriter, req *http.Request) {
	txn_addr, err := api.env.GetCurrentTransaction()
	if err != nil {
		panic("net/http/api: " + err.Error())
	}

	resp := struct {
		Transaction string
	}{
		Transaction: txn_addr.String(),
	}

	json.NewEncoder(w).Encode(resp)
}

func (api *API) handle_GET_view(w http.ResponseWriter, req *http.Request, handle *table_handle) {
	var (
		table = handle.table
		store = runtime.Env.Store
	)

	if table == nil {
		table = runtime.GetTable(store, handle.addr)
		handle.table = table
	}

	iter := table.Iter()

	w.Header().Set("Content-Type", "text/json; charset=utf-8")

	w.WriteHeader(200)
	w.Write([]byte("[\n"))

	first := true
	for {
		_, elt_addr, err := iter.Next()
		if err == btree.EOI {
			break
		}
		if err != nil {
			panic("net/http/api: " + err.Error())
		}

		var elt []byte

		err = cas.Decode(store, elt_addr, &elt)
		if err != nil {
			panic(fmt.Sprintf("net/http/api: (%+v) %s", elt_addr, err.Error()))
		}

		format := ",\n%s"
		if first {
			first = false
			format = "%s"
		}

		fmt.Fprintf(w, format, elt)
	}

	w.Write([]byte("\n]\n"))
}

func (api *API) handle_PATCH_transaction(w http.ResponseWriter, req *http.Request) {

}

package api

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/fd/simplex/cas"
	"github.com/fd/simplex/runtime"
	"net/http"
	"reflect"
	"strings"
)

type (
	API struct {
		name   string
		env    *runtime.Environment
		tables map[string]runtime.Table
		views  map[string]runtime.Deferred
		routes map[string]string

		ViewTables map[string]storage.SHA
	}
)

func New(env *runtime.Environment, name string) *API {
	api := &API{
		name,
		env,
		map[string]runtime.Table{},
		map[string]runtime.Deferred{},
		map[string]string{},
		map[string]storage.SHA{},
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

	json_view := runtime.Collect(
		view,
		func(ctx *runtime.Context, m_sha runtime.SHA) runtime.SHA {
			var m reflect.Value
			m = reflect.New(view.EltType())
			ctx.LoadValue(m_sha, m)

			data, err := json.Marshal(map[string]interface{}{
				"Vsn": hex.EncodeToString([]byte(m_sha[:])),
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

	for e := range funnel.Run() {
		event, ok := e.(*runtime.EvConsistent)
		if !ok {
			continue
		}

		if strings.HasPrefix(event.Table, "API/FORMAT_JSON/") {
			name := event.Table[len("API/FORMAT_JSON/"):]

			if event.B.IsZero() {
				delete(api.ViewTables, name)
			} else {
				api.ViewTables[name] = event.B
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
			if sha, p := api.ViewTables[view_id]; p {
				api.handle_GET_view(w, req, sha)
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
	txn_id, ok := api.env.GetCurrentTransaction()
	if !ok {
		panic("failed to get transaction id.")
	}

	resp := struct {
		TransactionID string
	}{
		TransactionID: hex.EncodeToString([]byte(txn_id[:])),
	}

	json.NewEncoder(w).Encode(resp)
}

func (api *API) handle_GET_view(w http.ResponseWriter, req *http.Request, sha storage.SHA) {
	var (
		table = runtime.Env.LoadTable(runtime.SHA(sha))
		store = runtime.Env.Store()
		iter  = table.Iter()
	)

	w.Header().Set("Content-Type", "text/json; charset=utf-8")

	w.WriteHeader(200)
	w.Write([]byte("[\n"))

	first := true
	for {
		sha, done := iter.Next()
		if done {
			break
		}

		var (
			kv  runtime.KeyValue
			val []byte
		)

		if !store.Get(sha, &kv) {
			panic("corrupt")
		}

		if !store.Get(kv.ValueSha, &val) {
			panic("corrupt")
		}

		format := ",\n%s"
		if first {
			first = false
			format = "%s"
		}

		fmt.Fprintf(w, format, val)
	}

	w.Write([]byte("\n]\n"))
}

func (api *API) handle_PATCH_transaction(w http.ResponseWriter, req *http.Request) {

}

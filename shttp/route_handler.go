package shttp

import (
	"encoding/json"
	"io"
	"net/http"
	"simplex.sh/store"
	"strings"
	"sync"
)

type RouteHandler struct {
	store store.Store
	hosts []string
	table map[string][]route_rule
	mtx   sync.RWMutex
}

func NewRouteHandler(store store.Store) (*RouteHandler, error) {
	m := &RouteHandler{store: store}

	err := m.load_routing_table()
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (r *RouteHandler) Hostnames() []string {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	return r.hosts
}

func (m *RouteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		rule = m.lookup(r)
		body io.ReadCloser
		err  error
	)

	if rule == nil {
		http.NotFound(w, r)
		return
	}

	if rule.Address != "" {
		body, err = m.store.GetBlob("blobs/" + rule.Address)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		defer body.Close()
	}

	{ // set the headers
		header := w.Header()
		for k, v := range rule.Header {
			header[k] = v
		}
	}

	w.WriteHeader(rule.Status)

	if body != nil {
		_, err := io.Copy(w, body)
		if err != nil {
			panic(err)
		}
	}
}

func (m *RouteHandler) load_routing_table() error {

	r, err := m.store.GetBlob("route_table.json")
	if err != nil {
		if store.IsNotFound(err) {
			m.table = map[string][]route_rule{}
			return nil
		} else {
			return err
		}
	}

	var (
		table    = map[string][]route_rule{}
		hosts    = []string{}
		host_map = map[string]bool{}
	)

	err = json.NewDecoder(r).Decode(&table)
	if err != nil {
		return err
	}

	for _, rules := range table {
		for _, rule := range rules {
			host := rule.Host

			if host == "." {
				continue
			}

			host_map[host] = true
		}
	}

	for host := range host_map {
		hosts = append(hosts, host)
	}

	m.mtx.Lock()
	defer m.mtx.Unlock()

	m.table = table
	m.hosts = hosts
	return nil
}

func (m *RouteHandler) lookup(r *http.Request) *route_rule {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	rules, p := m.table[r.URL.Path]
	if !p {
		return nil
	}

	var rule *route_rule

	for i := range rules {
		rule = &rules[i]

		if !rule.MatchingHost(r.Host) {
			continue
		}

		// match content-type
		// match language

		break
	}

	return rule
}

func (r *route_rule) MatchingHost(host string) bool {
	var (
		r_host   = r.Host
		wildcard bool
	)

	if strings.HasPrefix(r.Host, "*.") {
		r_host = r_host[2:]
		wildcard = true
	}

	if !strings.HasSuffix(host, ".") {
		host += "."
	}

	if !wildcard {
		return r_host == host
	}

	if !strings.HasSuffix(host, r_host) {
		return false
	}

	host = host[:len(host)-len(r_host)]
	return host == "" || strings.HasSuffix(host, ".")
}

package shttp

import (
	"net/http"
	"strings"
	"sync"
)

type HostHandler interface {
	http.Handler
	Hostnames() []string
}

type HostRouter struct {
	handlers []HostHandler
	table    map[string]HostHandler
	mtx      sync.RWMutex
}

func (h *HostRouter) Add(handler HostHandler) {
	h.handlers = append(h.handlers, handler)
}

func (h *HostRouter) Reset() {
	table := map[string]HostHandler{}

	for _, handler := range h.handlers {
		for _, hostname := range handler.Hostnames() {
			table[hostname] = handler
		}
	}

	h.set_table(table)
}

func (h *HostRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler := h.lookup(r)
	if handler == nil {
		http.NotFound(w, r)
		return
	}

	handler.ServeHTTP(w, r)
}

func (h *HostRouter) set_table(table map[string]HostHandler) {
	h.mtx.Lock()
	defer h.mtx.Unlock()

	h.table = table
}

func (h *HostRouter) get_table() map[string]HostHandler {
	h.mtx.RLock()
	defer h.mtx.RUnlock()

	return h.table
}

// www.example.com. matches:
// - www.example.com.
// - *.www.example.com.
// - *.example.com.
// - *.com.
// - *.
func (h *HostRouter) lookup(r *http.Request) http.Handler {
	table := h.get_table()

	host := r.Host

	if idx := strings.Index(host, ":"); idx >= 0 {
		host = host[:idx]
	}

	if !strings.HasSuffix(host, ".") {
		host += "."
	}

	if handler, p := table[host]; p {
		return handler
	}

	for {
		if handler, p := table["*."+host]; p {
			return handler
		}

		if host == "." {
			return nil
		}

		idx := strings.Index(host, ".")
		if idx < 0 {
			return nil
		}

		host = host[idx+1:]
	}

	panic("unreachable")
}

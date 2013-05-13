package shttp

import (
	"encoding/json"
	"net/http"
	"simplex.sh/errors"
	"sync"
)

type route_table_writer struct {
	Paths map[string]*route_set_writer
	mtx   sync.Mutex
}

type route_set_writer struct {
	Rules []*route_rule
	mtx   sync.Mutex
}

type route_rule struct {
	Host        string      `json:"h,omitempty"`
	Language    string      `json:"L,omitempty"`
	ContentType string      `json:"C,omitempty"`
	Status      int         `json:"S,omitempty"`
	Header      http.Header `json:"H,omitempty"`
	Address     string      `json:"A,omitempty"`
}

func (r *route_table_writer) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Paths)
}

func (r *route_table_writer) path(path string) *route_set_writer {
	if path == "" {
		path = "/"
	}

	r.mtx.Lock()
	defer r.mtx.Unlock()

	if r.Paths == nil {
		r.Paths = map[string]*route_set_writer{}
	}

	set := r.Paths[path]
	if set == nil {
		set = &route_set_writer{}
		r.Paths[path] = set
	}

	return set
}

func (r *route_set_writer) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Rules)
}

func (r *route_set_writer) add(rule *route_rule) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	for _, other := range r.Rules {
		if compare_lookup(other, rule) == 0 {

			if conflict(other, rule) {
				return errors.Fmt("Conflicting rule: %+v", rule)
			} else {
				return nil
			}

		}
	}

	r.Rules = append(r.Rules, rule)
	return nil
}

func compare_lookup(a, b *route_rule) int {
	if a.Host < b.Host {
		return -1
	}
	if a.Host > b.Host {
		return 1
	}

	if a.Language < b.Language {
		return -1
	}
	if a.Language > b.Language {
		return 1
	}

	if a.ContentType < b.ContentType {
		return -1
	}
	if a.ContentType > b.ContentType {
		return 1
	}

	return 0
}

func conflict(a, b *route_rule) bool {
	if a.Address != b.Address {
		return true
	}

	if a.Status != b.Status {
		return true
	}

	if len(a.Header) != len(b.Header) {
		return true
	}

	for k, a_h := range a.Header {
		b_h, p := b.Header[k]
		if !p {
			return true
		}

		if len(a_h) != len(b_h) {
			return true
		}

		for i, a_hv := range a_h {
			b_hv := b_h[i]
			if a_hv != b_hv {
				return true
			}
		}
	}

	return false
}

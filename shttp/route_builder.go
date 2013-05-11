package shttp

import (
	"strings"
)

type Router interface {
	Path(s string) Constrainer
}

type Constrainer interface {
	Host(s string)
}

type route_builder struct {
	rules []*route_rule_builder
}

type route_rule_builder struct {
	path string
	host string
}

func (r *route_builder) Path(s string) Constrainer {
	rule := &route_rule_builder{path: s}
	r.rules = append(r.rules, rule)
	return rule
}

func (r *route_rule_builder) Host(s string) {
	if !strings.HasSuffix(s, ".") {
		s += "."
	}

	r.host = s
}

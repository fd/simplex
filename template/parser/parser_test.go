package parser

import (
	"testing"
)

func TestParser(t *testing.T) {
	var (
		err error
	)

	tmpl, err := Parse("", `
    <a href="{{ person.link }}">{{ person.name }}</a>
    {{#if person.bio }}
      {{! This is a comment }}
      <a href={{ person.website().url }}>website</a>
      {{{ yield }}}
      {{ image_tag person.picture, title="Profile \"picture\"", class=person.greet("Hello") }}
    {{ else }}
      {{! No profile }}
    {{/if}}
  `)
	if err != nil {
		t.Errorf("error: %s", err)
	}
	e := `
    <a href="{{person.link}}">{{person.name}}</a>
    {{#if(person.bio)}}
      {{! This is a comment }}
      <a href={{person.website().url}}>website</a>
      {{{yield}}}
      {{image_tag(person.picture, class=person.greet("Hello"), title="Profile \"picture\"")}}
    {{else}}
      {{! No profile }}
    {{/end}}
  `
	if a := tmpl.String(); e != a {
		t.Errorf("\nexpected: %s\nactual: %s", e, a)
	}

}

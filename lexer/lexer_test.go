package lexer

import (
	"fmt"
	"strings"
	"testing"
)

func TestLexer(t *testing.T) {
	var (
		err error
		exp []Item
	)

	exp = []Item{
		{ItemEOF, ""},
	}
	err = consume(Lex("empty", ""), exp)
	if err != nil {
		t.Errorf("error: %s", err)
	}

	exp = []Item{
		{ItemHtmlText, "hello"},
		{ItemEOF, ""},
	}
	err = consume(Lex("", "hello"), exp)
	if err != nil {
		t.Errorf("error: %s", err)
	}

	exp = []Item{
		{ItemHtmlText, "Hello "},
		{ItemLeftMeta2, "{{"},
		{ItemIdentifier, "person.name"},
		{ItemRightMeta2, "}}"},
		{ItemHtmlText, "!"},
		{ItemEOF, ""},
	}
	err = consume(Lex("", "Hello {{ person.name }}!"), exp)
	if err != nil {
		t.Errorf("error: %s", err)
	}

	exp = []Item{
		{ItemHtmlLiteral, "<a href"},
		{ItemHtmlAttrInterp, "=\""},
		{ItemLeftMeta2, "{{"},
		{ItemIdentifier, "person.link"},
		{ItemRightMeta2, "}}"},
		{ItemHtmlLiteral, "\">"},
		{ItemLeftMeta2, "{{"},
		{ItemIdentifier, "person.name"},
		{ItemRightMeta2, "}}"},
		{ItemHtmlLiteral, "</a>"},
		{ItemHtmlText, "\n    "},
		{ItemLeftMetaPound, "{{#"},
		{ItemIdentifier, "if"},
		{ItemIdentifier, "person.bio"},
		{ItemRightMeta2, "}}"},
		{ItemHtmlText, "\n      "},
		{ItemLeftMetaBang, "{{!"},
		{ItemHtmlComment, " This is a comment "},
		{ItemRightMeta2, "}}"},
		{ItemHtmlText, "\n      "},
		{ItemHtmlLiteral, "<a href"},
		{ItemHtmlAttr, "={{"},
		{ItemIdentifier, "person.website"},
		{ItemRightMeta2, "}}"},
		{ItemHtmlLiteral, ">"},
		{ItemHtmlText, "website"},
		{ItemHtmlLiteral, "</a>"},
		{ItemHtmlText, "\n      "},
		{ItemLeftMeta3, "{{{"},
		{ItemIdentifier, "yield"},
		{ItemRightMeta3, "}}}"},
		{ItemHtmlText, "\n      "},
		{ItemLeftMeta2, "{{"},
		{ItemIdentifier, "image_tag"},
		{ItemIdentifier, "person.picture"},
		{ItemIdentifier, "title"},
		{ItemEqual, "="},
		{ItemString, "Profile \\\"picture\\\""},
		{ItemIdentifier, "class"},
		{ItemEqual, "="},
		{ItemIdentifier, "person.hello"},
		{ItemRightMeta2, "}}"},
		{ItemHtmlText, "\n    "},
		{ItemLeftMetaSlash, "{{/"},
		{ItemIdentifier, "end"},
		{ItemRightMeta2, "}}"},
		{ItemEOF, ""},
	}
	err = consume(Lex("", strings.TrimSpace(`
    <a href="{{ person.link }}">{{ person.name }}</a>
    {{#if person.bio }}
      {{! This is a comment }}
      <a href={{ person.website }}>website</a>
      {{{ yield }}}
      {{ image_tag person.picture title="Profile \"picture\"" class=person.hello }}
    {{/end}}
  `)), exp)
	if err != nil {
		t.Errorf("error: %s", err)
	}
}

func consume(c <-chan Item, expected []Item) error {
	actual := []Item{}

	for item := range c {
		actual = append(actual, item)
	}

	for i, e := range expected {
		if i >= len(actual) {
			return fmt.Errorf("missing tokens: %+v", expected[i:])
		}

		if actual[i] != e {
			return fmt.Errorf("Unexpected token %+v (expected %+v)", actual[i], e)
		}
	}

	if len(actual) > len(expected) {
		return fmt.Errorf("unexpected tokens: %+v", actual[len(expected):])
	}

	return nil
}

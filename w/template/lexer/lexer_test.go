package lexer

import (
	"fmt"
	"testing"
)

func TestLexer(t *testing.T) {
	var (
		err error
		exp []Item
	)

	exp = []Item{
		{ItemEOF, "", 1, 1},
	}
	err = consume(Lex("empty", ""), exp)
	if err != nil {
		t.Errorf("error: %s", err)
	}

	exp = []Item{
		{ItemHtmlText, "hello", 1, 1},
		{ItemEOF, "", 1, 6},
	}
	err = consume(Lex("", "hello"), exp)
	if err != nil {
		t.Errorf("error: %s", err)
	}

	exp = []Item{
		{ItemHtmlText, "Hello ", 1, 1},
		{ItemLeftMeta2, "{{", 1, 7},
		{ItemExpression, " person.name ", 1, 9},
		{ItemRightMeta2, "}}", 1, 22},
		{ItemHtmlText, "!", 1, 24},
		{ItemEOF, "", 1, 25},
	}
	err = consume(Lex("", "Hello {{ person.name }}!"), exp)
	if err != nil {
		t.Errorf("error: %s", err)
	}

	exp = []Item{
		{ItemHtmlText, "\n    ", 1, 1},
		{ItemHtmlLiteral, "<a href", 2, 5},
		{ItemHtmlAttrInterp, "=\"", 2, 12},
		{ItemLeftMeta2, "{{", 2, 14},
		{ItemExpression, " person.link ", 2, 16},
		{ItemRightMeta2, "}}", 2, 29},
		{ItemHtmlAttrInterpEnd, "\"", 2, 31},
		{ItemHtmlLiteral, ">", 2, 32},
		{ItemLeftMeta2, "{{", 2, 33},
		{ItemExpression, " person.name ", 2, 35},
		{ItemRightMeta2, "}}", 2, 48},
		{ItemHtmlLiteral, "</a>", 2, 50},
		{ItemHtmlText, "\n    ", 2, 54},
		{ItemLeftMetaPound, "{{#", 3, 5},
		{ItemIdentifier, "if", 3, 8},
		{ItemExpression, " person.bio ", 3, 10},
		{ItemRightMeta2, "}}", 3, 22},
		{ItemHtmlText, "\n      ", 3, 24},
		{ItemLeftMetaBang, "{{!", 4, 7},
		{ItemHtmlComment, " This is a comment ", 4, 10},
		{ItemRightMeta2, "}}", 4, 29},
		{ItemHtmlText, "\n      ", 4, 31},
		{ItemHtmlLiteral, "<a href", 5, 7},
		{ItemHtmlAttr, "={{", 5, 14},
		{ItemExpression, " person.website ", 5, 17},
		{ItemRightMeta2, "}}", 5, 33},
		{ItemHtmlLiteral, ">", 5, 35},
		{ItemHtmlText, "website", 5, 36},
		{ItemHtmlLiteral, "</a>", 5, 43},
		{ItemHtmlText, "\n      ", 5, 47},
		{ItemLeftMeta3, "{{{", 6, 7},
		{ItemExpression, " yield ", 6, 10},
		{ItemRightMeta3, "}}}", 6, 17},
		{ItemHtmlText, "\n      ", 6, 20},
		{ItemLeftMeta2, "{{", 7, 7},
		{ItemExpression, " image_tag( person.picture, title=\"Profile \\\"picture\\\"\", class=person.hello ) ", 7, 9},
		{ItemRightMeta2, "}}", 7, 87},
		{ItemHtmlText, "\n    ", 7, 89},
		{ItemLeftMetaSlash, "{{/", 8, 5},
		{ItemIdentifier, "end", 8, 8},
		{ItemRightMeta2, "}}", 8, 11},
		{ItemHtmlText, "\n  ", 8, 13},
		{ItemEOF, "", 9, 3},
	}
	err = consume(Lex("", `
    <a href="{{ person.link }}">{{ person.name }}</a>
    {{#if person.bio }}
      {{! This is a comment }}
      <a href={{ person.website }}>website</a>
      {{{ yield }}}
      {{ image_tag( person.picture, title="Profile \"picture\"", class=person.hello ) }}
    {{/end}}
  `), exp)
	if err != nil {
		t.Errorf("error: %s", err)
	}

	exp = []Item{
		{ItemHtmlLiteral, "<a ", 1, 1},
		{ItemError, "unexpected {{ token", 1, 4},
	}
	err = consume(Lex("", "<a {{ person }}>"), exp)
	if err != nil {
		t.Errorf("error: %s", err)
	}

	exp = []Item{
		{ItemHtmlLiteral, "<a ", 1, 1},
		{ItemError, "unexpected {{ token", 1, 4},
	}
	err = consume(Lex("", "<a {{{ person }}}>"), exp)
	if err != nil {
		t.Errorf("error: %s", err)
	}

	exp = []Item{
		{ItemHtmlLiteral, "<a href", 1, 1},
		{ItemHtmlAttr, "={{", 1, 8},
		{ItemExpression, "\n  person\n    .name\n", 1, 11},
		{ItemRightMeta2, "}}", 4, 1},
		{ItemHtmlLiteral, ">", 4, 3},
		{ItemEOF, "", 4, 4},
	}
	err = consume(Lex("", "<a href={{\n  person\n    .name\n}}>"), exp)
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

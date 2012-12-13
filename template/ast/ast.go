package ast

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type Context int

const (
	TextContext Context = iota
	AttrContext
	AttrTextContext
)

type Node interface {
	NodeInfo() Info
	String() string
}

type Statement interface {
	Node
}

type Expression interface {
	Node
}

type Info struct {
	File   string
	Line   int
	Column int
}

func (i Info) NodeInfo() Info {
	return i
}

type Template struct {
	Info
	Statements []Statement
}

func (t *Template) String() string {
	s := ""
	for _, stmt := range t.Statements {
		s += stmt.String()
	}
	return s
}

type Block struct {
	Info
	Expression   Expression
	Template     *Template
	ElseTemplate *Template
}

func (b *Block) String() string {
	if b.ElseTemplate != nil && len(b.ElseTemplate.Statements) > 0 {
		return fmt.Sprintf("{{#%s}}%s{{else}}%s{{/end}}", b.Expression, b.Template, b.ElseTemplate)
	}
	return fmt.Sprintf("{{#%s}}%s{{/end}}", b.Expression, b.Template)
}

type Interpolation struct {
	Info
	Expression Expression
	Raw        bool
	Context    Context
}

func (i *Interpolation) String() string {
	if i.Raw {
		return fmt.Sprintf("{{{%s}}}", i.Expression)
	}
	return fmt.Sprintf("{{%s}}", i.Expression)
}

type Comment struct {
	Info
	Content string
}

func (c *Comment) String() string {
	return fmt.Sprintf("{{!%s}}", c.Content)
}

type Literal struct {
	Info
	Content string
}

func (l *Literal) String() string {
	return l.Content
}

type IntegerLiteral struct {
	Info
	Value int
}

func (n *IntegerLiteral) String() string {
	return fmt.Sprintf("%d", n.Value)
}

type FloatLiteral struct {
	Info
	Value float64
}

func (n *FloatLiteral) String() string {
	return fmt.Sprintf("%f", n.Value)
}

type StringLiteral struct {
	Info
	Value string
}

func (n *StringLiteral) String() string {
	return strconv.Quote(n.Value)
}

type Identifier struct {
	Info
	Value string
}

func (n *Identifier) String() string {
	return n.Value
}

type Get struct {
	Info
	From Expression
	Name *Identifier
}

func (n *Get) String() string {
	return n.From.String() + "." + n.Name.String()
}

type FunctionCall struct {
	Info
	From    Expression
	Name    string
	Args    []Expression
	Options map[string]Expression
}

func (n *FunctionCall) String() string {
	from := ""
	if n.From != nil {
		from = n.From.String() + "."
	}

	args := make([]string, 0, len(n.Args)+len(n.Options))

	for _, arg := range n.Args {
		args = append(args, arg.String())
	}

	keys := make([]string, 0, len(n.Options))
	for key := range n.Options {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		args = append(args, key+"="+n.Options[key].String())
	}

	return from + n.Name + "(" + strings.Join(args, ", ") + ")"
}

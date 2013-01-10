package ast

import (
	"fmt"
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

func (t *Template) GoString() string {
	s := ""
	for _, stmt := range t.Statements {
		s += stmt.String()
	}
	return s
}

type Macro struct {
	Info
	Macro      *Identifier
	Expression *Expression
}

func (n *Macro) String() string {
	return fmt.Sprintf("{{:%s%s}}", n.Macro, n.Expression)
}

type Block struct {
	Info
	Macro        *Identifier
	Expression   *Expression
	Template     *Template
	ElseTemplate *Template
}

func (b *Block) String() string {
	if b.ElseTemplate != nil && len(b.ElseTemplate.Statements) > 0 {
		return fmt.Sprintf("{{#%s%s}}%s{{else}}%s{{/end}}", b.Macro, b.Expression, b.Template, b.ElseTemplate)
	}
	return fmt.Sprintf("{{#%s%s}}%s{{/end}}", b.Macro, b.Expression, b.Template)
}

type Interpolation struct {
	Info
	Expression *Expression
	Raw        bool
	Context    Context
}

func (i *Interpolation) String() string {
	if i.Raw {
		return fmt.Sprintf("{{{%s}}}", i.Expression)
	}
	return fmt.Sprintf("{{%s}}", i.Expression)
}

type Literal struct {
	Info
	Content string
}

func (l *Literal) String() string {
	return l.Content
}

type Comment struct {
	Info
	Content string
}

func (c *Comment) String() string {
	return fmt.Sprintf("{{!%s}}", c.Content)
}

type Identifier struct {
	Info
	Value string
}

func (n *Identifier) String() string {
	return n.Value
}

type Expression struct {
	Info
	Value string
}

func (n *Expression) String() string {
	return n.Value
}

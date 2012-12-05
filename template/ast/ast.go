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
	Visit(Visitor)
	FirstNode() Node
	NodeInfo() Info
}

type Statement interface {
	Node
	String() string
}

type Expression interface {
	Node
	String() string
}

type Visitor interface {
	VisitTemplate(*Template)
	VisitBlock(*Block)
	VisitComment(*Comment)
	VisitLiteral(*Literal)
	VisitInterpolation(*Interpolation)

	VisitStringLiteral(*StringLiteral)
	VisitIntegerLiteral(*IntegerLiteral)
	VisitFloatLiteral(*FloatLiteral)
	VisitIdentifier(*Identifier)
	VisitGet(*Get)
	VisitFunctionCall(*FunctionCall)
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

func (n *Template) Visit(v Visitor) {
	v.VisitTemplate(n)

	for _, stmt := range n.Statements {
		stmt.Visit(v)
	}
}

func (n *Template) FirstNode() Node {
	if len(n.Statements) > 0 {
		return n.Statements[0].FirstNode()
	}
	return n
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

func (b *Block) Visit(v Visitor) {
	v.VisitBlock(b)

	b.Expression.Visit(v)
	b.Template.Visit(v)
	b.ElseTemplate.Visit(v)
}

func (n *Block) FirstNode() Node {
	return n.Expression.FirstNode()
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

func (n *Interpolation) Visit(v Visitor) {
	v.VisitInterpolation(n)

	n.Expression.Visit(v)
}

func (n *Interpolation) FirstNode() Node {
	return n.Expression.FirstNode()
}

type Comment struct {
	Info
	Content string
}

func (c *Comment) String() string {
	return fmt.Sprintf("{{!%s}}", c.Content)
}

func (c *Comment) Visit(v Visitor) {
	v.VisitComment(c)
}

func (n *Comment) FirstNode() Node {
	return n
}

type Literal struct {
	Info
	Content string
}

func (l *Literal) String() string {
	return l.Content
}

func (l *Literal) Visit(v Visitor) {
	v.VisitLiteral(l)
}

func (n *Literal) FirstNode() Node {
	return n
}

type IntegerLiteral struct {
	Info
	Value int
}

func (n *IntegerLiteral) String() string {
	return fmt.Sprintf("%d", n.Value)
}

func (n *IntegerLiteral) Visit(v Visitor) {
	v.VisitIntegerLiteral(n)
}

func (n *IntegerLiteral) FirstNode() Node {
	return n
}

type FloatLiteral struct {
	Info
	Value float64
}

func (n *FloatLiteral) String() string {
	return fmt.Sprintf("%f", n.Value)
}

func (n *FloatLiteral) Visit(v Visitor) {
	v.VisitFloatLiteral(n)
}

func (n *FloatLiteral) FirstNode() Node {
	return n
}

type StringLiteral struct {
	Info
	Value string
}

func (n *StringLiteral) Visit(v Visitor) {
	v.VisitStringLiteral(n)
}

func (n *StringLiteral) String() string {
	return strconv.Quote(n.Value)
}

func (n *StringLiteral) FirstNode() Node {
	return n
}

type Identifier struct {
	Info
	Value string
}

func (n *Identifier) String() string {
	return n.Value
}

func (n *Identifier) Visit(v Visitor) {
	v.VisitIdentifier(n)
}

func (n *Identifier) FirstNode() Node {
	return n
}

type Get struct {
	Info
	From Expression
	Name *Identifier
}

func (n *Get) String() string {
	return n.From.String() + "." + n.Name.String()
}

func (n *Get) Visit(v Visitor) {
	v.VisitGet(n)

	n.From.Visit(v)
	n.Name.Visit(v)
}

func (n *Get) FirstNode() Node {
	return n.From.FirstNode()
}

type FunctionCall struct {
	Info
	From    Expression
	Args    []Expression
	Options map[string]Expression
}

func (n *FunctionCall) String() string {
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

	return n.From.String() + "(" + strings.Join(args, ", ") + ")"
}

func (n *FunctionCall) Visit(v Visitor) {
	v.VisitFunctionCall(n)

	n.From.Visit(v)

	for _, arg := range n.Args {
		arg.Visit(v)
	}

	for _, arg := range n.Options {
		arg.Visit(v)
	}
}

func (n *FunctionCall) FirstNode() Node {
	return n.From.FirstNode()
}

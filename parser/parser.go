package parser

import (
	"fmt"
	"github.com/fd/handlebars/lexer"
	"strconv"
	"strings"
)

type Context int

const (
	TextContext Context = iota
	AttrContext
	AttrTextContext
)

type Template []Statement

func (t Template) String() string {
	s := ""
	for _, stmt := range t {
		s += stmt.String()
	}
	return s
}

type Statement interface {
	Visit(Visitor)
	String() string
}

type Expression interface {
	String() string
}

type Visitor interface {
	EnterTemplate()
	LeaveTemplate()

	EnterBlock(*Block)
	LeaveBlock(*Block)

	VisitComment(*Comment)
	VisitHelper()
	VisitValue(*Interpolation)
	VisitLiteral(*Literal)
}

type Block struct {
	Expression   Expression
	Template     Template
	ElseTemplate Template
}

func (b *Block) String() string {
	if len(b.ElseTemplate) > 0 {
		return fmt.Sprintf("{{#%s}}%s{{else}}%s{{/end}}", b.Expression, b.Template, b.ElseTemplate)
	}
	return fmt.Sprintf("{{#%s}}%s{{/end}}", b.Expression, b.Template)
}

func (b *Block) Visit(v Visitor) {
	v.EnterBlock(b)
	// for
	v.LeaveBlock(b)
}

type Interpolation struct {
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

func (i *Interpolation) Visit(v Visitor) {
	v.VisitValue(i)
}

type Comment struct {
	Content string
}

func (c *Comment) String() string {
	return fmt.Sprintf("{{!%s}}", c.Content)
}

func (c *Comment) Visit(v Visitor) {
	v.VisitComment(c)
}

type Literal struct {
	Content string
}

func (l *Literal) String() string {
	return l.Content
}

func (l *Literal) Visit(v Visitor) {
	v.VisitLiteral(l)
}

type IntegerLiteral struct {
	Value int
}

func (n *IntegerLiteral) String() string {
	return fmt.Sprintf("%d", n.Value)
}

type FloatLiteral struct {
	Value float64
}

func (n *FloatLiteral) String() string {
	return fmt.Sprintf("%f", n.Value)
}

type StringLiteral struct {
	Value string
}

func (n *StringLiteral) String() string {
	return strconv.Quote(n.Value)
}

type Identifier struct {
	Value string
}

func (n *Identifier) String() string {
	return n.Value
}

type Get struct {
	From Expression
	Name *Identifier
}

func (n *Get) String() string {
	return n.From.String() + "." + n.Name.String()
}

type FunctionCall struct {
	From    Expression
	Args    []Expression
	Options map[string]Expression
}

func (n *FunctionCall) String() string {
	args := make([]string, 0, len(n.Args)+len(n.Options))

	for _, arg := range n.Args {
		args = append(args, arg.String())
	}

	for key, arg := range n.Options {
		args = append(args, key+"="+arg.String())
	}

	return n.From.String() + "(" + strings.Join(args, ", ") + ")"
}

func Parse(name, content string) (Template, error) {
	p := &parser{tokenChan: lexer.Lex(name, content)}
	return p.parseTemplate()
}

type parser struct {
	tokenChan <-chan lexer.Item
	backlog   []lexer.Item
	level     int
}

func (p *parser) pop() lexer.Item {
	if len(p.backlog) > 0 {
		token := p.backlog[len(p.backlog)-1]
		p.backlog = p.backlog[:len(p.backlog)-1]
		return token
	}

	token, ok := <-p.tokenChan
	if !ok {
		return lexer.Item{lexer.ItemEOF, "", 0, 0}
	}

	return token
}

func (p *parser) push(token lexer.Item) {
	if p.backlog == nil {
		p.backlog = make([]lexer.Item, 0, 100)
	}
	p.backlog = append(p.backlog, token)
}

func (p *parser) parseTemplate() (Template, error) {
	tmpl := Template{}
	ctx := TextContext

	for {
		token := p.pop()

		switch token.Type {

		case lexer.ItemEOF:
			return tmpl, nil

		case lexer.ItemLeftMetaSlash:
			p.push(token)
			return tmpl, nil

		case lexer.ItemHtmlText:
			tmpl = append(tmpl, &Literal{token.Value})

		case lexer.ItemHtmlLiteral:
			tmpl = append(tmpl, &Literal{token.Value})

		case lexer.ItemHtmlAttr:
			p.push(token)
			stmt, err := p.parseHtmlAttr()
			if err != nil {
				return nil, err
			}
			tmpl = append(tmpl, &Literal{"="})
			tmpl = append(tmpl, stmt)

		case lexer.ItemHtmlAttrInterp:
			tmpl = append(tmpl, &Literal{token.Value})
			ctx = AttrTextContext

		case lexer.ItemHtmlAttrInterpEnd:
			tmpl = append(tmpl, &Literal{token.Value})
			ctx = TextContext

		case lexer.ItemLeftMeta2, lexer.ItemLeftMeta3:
			{ // look for {{else}}
				if token.Type == lexer.ItemLeftMeta2 {
					n_else := p.pop()
					n_close := p.pop()

					if n_else.Type == lexer.ItemIdentifier && n_else.Value == "else" && n_close.Type == lexer.ItemRightMeta2 {
						p.push(n_close)
						p.push(n_else)
						p.push(token)
						return tmpl, nil
					} else {
						p.push(n_close)
						p.push(n_else)
					}
				}
			}

			p.push(token)
			stmt, err := p.parseInterpolation(ctx)
			if err != nil {
				return nil, err
			}
			tmpl = append(tmpl, stmt)

		case lexer.ItemLeftMetaBang:
			p.push(token)
			stmt, err := p.parseComment()
			if err != nil {
				return nil, err
			}
			tmpl = append(tmpl, stmt)

		case lexer.ItemLeftMetaPound:
			p.push(token)
			stmt, err := p.parseBlock()
			if err != nil {
				return nil, err
			}
			tmpl = append(tmpl, stmt)

		default:
			return nil, p.unexpected_token(token)

		}
	}

	return nil, fmt.Errorf("Unexpected EOF")
}

func (p *parser) parseHtmlAttr() (Statement, error) {
	if token := p.pop(); token.Type != lexer.ItemHtmlAttr {
		return nil, p.unexpected_token(token)
	}

	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if token := p.pop(); token.Type != lexer.ItemRightMeta2 {
		return nil, p.unexpected_token(token)
	}

	return &Interpolation{expr, false, AttrContext}, nil
}

func (p *parser) parseComment() (Statement, error) {
	if token := p.pop(); token.Type != lexer.ItemLeftMetaBang {
		return nil, p.unexpected_token(token)
	}

	comment := &Comment{Content: ""}
	if token := p.pop(); token.Type != lexer.ItemHtmlComment {
		return nil, p.unexpected_token(token)
	} else {
		comment.Content = token.Value
	}

	if token := p.pop(); token.Type != lexer.ItemRightMeta2 {
		return nil, p.unexpected_token(token)
	}

	return comment, nil
}

func (p *parser) parseInterpolation(ctx Context) (Statement, error) {
	var token lexer.Item

	token = p.pop()
	raw := false

	switch token.Type {
	case lexer.ItemLeftMeta2:
		raw = false
	case lexer.ItemLeftMeta3:
		raw = true
	default:
		return nil, p.unexpected_token(token)
	}

	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if raw {
		if token := p.pop(); token.Type != lexer.ItemRightMeta3 {
			return nil, p.unexpected_token(token)
		}
	} else {
		if token := p.pop(); token.Type != lexer.ItemRightMeta2 {
			return nil, p.unexpected_token(token)
		}
	}

	return &Interpolation{expr, raw, ctx}, nil
}

func (p *parser) parseBlock() (Statement, error) {
	var token lexer.Item

	if t := p.pop(); t.Type != lexer.ItemLeftMetaPound {
		return nil, p.unexpected_token(t)
	}

	token = p.pop()
	if token.Type != lexer.ItemIdentifier {
		return nil, p.unexpected_token(token)
	}
	name := token.Value
	p.push(token)

	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if t := p.pop(); t.Type != lexer.ItemRightMeta2 {
		return nil, p.unexpected_token(t)
	}

	var t_else Template
	t_main, err := p.parseTemplate()
	if err != nil {
		return nil, err
	}

	if t := p.pop(); t.Type == lexer.ItemLeftMeta2 {
		n_else := p.pop()
		n_close := p.pop()

		if n_else.Type == lexer.ItemIdentifier && n_else.Value == "else" && n_close.Type == lexer.ItemRightMeta2 {
			t_else, err = p.parseTemplate()
			if err != nil {
				return nil, err
			}
		} else {
			p.push(n_close)
			p.push(n_else)
			p.push(t)
		}
	} else {
		p.push(t)
	}

	if t := p.pop(); t.Type != lexer.ItemLeftMetaSlash {
		return nil, p.unexpected_token(t)
	}

	token = p.pop()
	if token.Type != lexer.ItemIdentifier {
		return nil, p.unexpected_token(token)
	}
	if t := p.pop(); t.Type != lexer.ItemRightMeta2 {
		return nil, p.unexpected_token(t)
	}
	if token.Value != "end" && token.Value != name {
		return nil, fmt.Errorf("Unmatched block close tag: {{/%s}}", token.Value)
	}

	return &Block{Expression: expr, Template: t_main, ElseTemplate: t_else}, nil
}

func (p *parser) parseExpression() (Expression, error) {
	t := p.pop()

	switch t.Type {
	case lexer.ItemNumber:
		p.push(t)
		return p.parseNumberExpression()

	case lexer.ItemString:
		p.push(t)
		return p.parseStringExpression()

	case lexer.ItemIdentifier:
		p.push(t)
		return p.parseIdentifierExpression()

	}

	return nil, p.unexpected_token(t)
}

func (p *parser) parseNumberExpression() (Expression, error) {
	t := p.pop()

	if i, err := strconv.Atoi(t.Value); err == nil {
		return &IntegerLiteral{Value: i}, nil
	}

	if f, err := strconv.ParseFloat(t.Value, 64); err == nil {
		return &FloatLiteral{Value: f}, nil
	}

	return nil, fmt.Errorf("Invalid number literal: %s", t.Value)
}

func (p *parser) parseStringExpression() (Expression, error) {
	t := p.pop()

	if s, err := strconv.Unquote(t.Value); err == nil {
		return &StringLiteral{Value: s}, nil
	}

	return nil, fmt.Errorf("Invalid string literal: %s", t.Value)
}

func (p *parser) parseAfterExpression(expr Expression) (Expression, error) {
	t := p.pop()

	if p.level > 1 {
		switch t.Type {
		case lexer.ItemDot:
			return p.parseGetExpression(expr)

		case lexer.ItemLeftParen:
			return p.parseFunctionCallExpression(expr, false)

		}

	} else {

		switch t.Type {
		case lexer.ItemDot:
			return p.parseGetExpression(expr)

		case lexer.ItemLeftParen, lexer.ItemString, lexer.ItemNumber, lexer.ItemIdentifier:
			bare := t.Type != lexer.ItemLeftParen
			if bare {
				p.push(t)
			}
			return p.parseFunctionCallExpression(expr, bare)

		}
	}

	p.push(t)

	return expr, nil
}

func (p *parser) parseIdentifierExpression() (Expression, error) {
	t := p.pop()
	expr := &Identifier{Value: t.Value}
	return p.parseAfterExpression(expr)
}

func (p *parser) parseGetExpression(base Expression) (Expression, error) {
	t := p.pop()
	if t.Type != lexer.ItemIdentifier {
		return nil, p.unexpected_token(t)
	}

	expr := &Get{From: base, Name: &Identifier{Value: t.Value}}
	return p.parseAfterExpression(expr)
}

func (p *parser) parseFunctionCallExpression(base Expression, bare bool) (Expression, error) {
	first := true
	args := []Expression{}
	opts := map[string]Expression{}

	p.level += 1

	for {
		t1 := p.pop()

		if !bare {
			if t1.Type == lexer.ItemRightParen {
				break
			}
		} else {
			if t1.Type == lexer.ItemRightMeta2 || t1.Type == lexer.ItemRightMeta3 {
				p.push(t1)
				break
			}
		}

		if !first {
			if t1.Type != lexer.ItemComma {
				return nil, p.unexpected_token(t1)
			}

			t1 = p.pop()
		}
		first = false

		regular := true

		if t1.Type == lexer.ItemIdentifier {
			t2 := p.pop()

			if t2.Type == lexer.ItemEqual {
				expr, err := p.parseExpression()
				if err != nil {
					return nil, err
				}
				opts[t1.Value] = expr
				regular = false
			} else {
				p.push(t2)
				p.push(t1)
			}
		} else {
			p.push(t1)
		}

		if regular {
			expr, err := p.parseExpression()
			if err != nil {
				return nil, err
			}
			args = append(args, expr)
		}
	}

	expr := &FunctionCall{From: base, Args: args, Options: opts}

	p.level -= 1
	return p.parseAfterExpression(expr)
}

func (p *parser) unexpected_token(token lexer.Item) error {
	panic(fmt.Errorf("%d:%d unexpected token: '%s'", token.Line, token.Column, token.Value))
	return fmt.Errorf("%d:%d unexpected token: '%s'", token.Line, token.Column, token.Value)
}

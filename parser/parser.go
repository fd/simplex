package parser

import (
	"fmt"
	"github.com/fd/handlebars/ast"
	"github.com/fd/handlebars/lexer"
	"strconv"
)

func Parse(name, content string) (*ast.Template, error) {
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

func (p *parser) parseTemplate() (*ast.Template, error) {
	tmpl := &ast.Template{}
	ctx := ast.TextContext

	for {
		token := p.pop()
		info := ast.Info{Line: token.Line, Column: token.Column}

		switch token.Type {

		case lexer.ItemEOF:
			return tmpl, nil

		case lexer.ItemLeftMetaSlash:
			p.push(token)
			return tmpl, nil

		case lexer.ItemHtmlText:
			tmpl.Statements = append(tmpl.Statements, &ast.Literal{info, token.Value})

		case lexer.ItemHtmlLiteral:
			tmpl.Statements = append(tmpl.Statements, &ast.Literal{info, token.Value})

		case lexer.ItemHtmlAttr:
			p.push(token)
			stmt, err := p.parseHtmlAttr()
			if err != nil {
				return nil, err
			}
			tmpl.Statements = append(tmpl.Statements, &ast.Literal{info, "="})
			tmpl.Statements = append(tmpl.Statements, stmt)

		case lexer.ItemHtmlAttrInterp:
			tmpl.Statements = append(tmpl.Statements, &ast.Literal{info, token.Value})
			ctx = ast.AttrTextContext

		case lexer.ItemHtmlAttrInterpEnd:
			tmpl.Statements = append(tmpl.Statements, &ast.Literal{info, token.Value})
			ctx = ast.TextContext

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
			tmpl.Statements = append(tmpl.Statements, stmt)

		case lexer.ItemLeftMetaBang:
			p.push(token)
			stmt, err := p.parseComment()
			if err != nil {
				return nil, err
			}
			tmpl.Statements = append(tmpl.Statements, stmt)

		case lexer.ItemLeftMetaPound:
			p.push(token)
			stmt, err := p.parseBlock()
			if err != nil {
				return nil, err
			}
			tmpl.Statements = append(tmpl.Statements, stmt)

		default:
			return nil, p.unexpected_token(token)

		}
	}

	return nil, fmt.Errorf("Unexpected EOF")
}

func (p *parser) parseHtmlAttr() (ast.Statement, error) {
	token := p.pop()
	if token.Type != lexer.ItemHtmlAttr {
		return nil, p.unexpected_token(token)
	}
	// ={{
	info := ast.Info{Line: token.Line, Column: token.Column + 1}

	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if token := p.pop(); token.Type != lexer.ItemRightMeta2 {
		return nil, p.unexpected_token(token)
	}

	return &ast.Interpolation{info, expr, false, ast.AttrContext}, nil
}

func (p *parser) parseComment() (ast.Statement, error) {
	token := p.pop()
	if token.Type != lexer.ItemLeftMetaBang {
		return nil, p.unexpected_token(token)
	}
	info := ast.Info{Line: token.Line, Column: token.Column + 1}

	comment := &ast.Comment{Info: info, Content: ""}
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

func (p *parser) parseInterpolation(ctx ast.Context) (ast.Statement, error) {
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
	info := ast.Info{Line: token.Line, Column: token.Column}

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

	return &ast.Interpolation{info, expr, raw, ctx}, nil
}

func (p *parser) parseBlock() (ast.Statement, error) {
	var token lexer.Item

	token = p.pop()
	if token.Type != lexer.ItemLeftMetaPound {
		return nil, p.unexpected_token(token)
	}
	info := ast.Info{Line: token.Line, Column: token.Column}

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

	var t_else *ast.Template
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

	return &ast.Block{info, expr, t_main, t_else}, nil
}

func (p *parser) parseExpression() (ast.Expression, error) {
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

func (p *parser) parseNumberExpression() (ast.Expression, error) {
	t := p.pop()
	info := ast.Info{Line: t.Line, Column: t.Column}

	if i, err := strconv.Atoi(t.Value); err == nil {
		return &ast.IntegerLiteral{info, i}, nil
	}

	if f, err := strconv.ParseFloat(t.Value, 64); err == nil {
		return &ast.FloatLiteral{info, f}, nil
	}

	return nil, fmt.Errorf("Invalid number literal: %s", t.Value)
}

func (p *parser) parseStringExpression() (ast.Expression, error) {
	t := p.pop()
	info := ast.Info{Line: t.Line, Column: t.Column}

	if s, err := strconv.Unquote(t.Value); err == nil {
		return &ast.StringLiteral{info, s}, nil
	}

	return nil, fmt.Errorf("Invalid string literal: %s", t.Value)
}

func (p *parser) parseAfterExpression(expr ast.Expression) (ast.Expression, error) {
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

func (p *parser) parseIdentifierExpression() (ast.Expression, error) {
	t := p.pop()
	expr := &ast.Identifier{Value: t.Value}
	return p.parseAfterExpression(expr)
}

func (p *parser) parseGetExpression(base ast.Expression) (ast.Expression, error) {
	t := p.pop()
	if t.Type != lexer.ItemIdentifier {
		return nil, p.unexpected_token(t)
	}

	info := ast.Info{Line: t.Line, Column: t.Column}
	expr := &ast.Get{info, base, &ast.Identifier{Value: t.Value}}
	return p.parseAfterExpression(expr)
}

func (p *parser) parseFunctionCallExpression(base ast.Expression, bare bool) (ast.Expression, error) {
	first := true
	args := []ast.Expression{}
	opts := map[string]ast.Expression{}

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

	info := base.FirstNode().NodeInfo()
	expr := &ast.FunctionCall{Info: info, From: base, Args: args, Options: opts}

	p.level -= 1
	return p.parseAfterExpression(expr)
}

func (p *parser) unexpected_token(token lexer.Item) error {
	return fmt.Errorf("%d:%d unexpected token: '%s'", token.Line, token.Column, token.Value)
}

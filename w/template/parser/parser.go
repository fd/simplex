package parser

import (
	"fmt"
	"io/ioutil"
	"simplex.sh/w/template/ast"
	"simplex.sh/w/template/lexer"
)

func ParseFile(path string) (*ast.Template, error) {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Parse(path, string(dat))
}

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

		case lexer.ItemLeftMetaColon:
			p.push(token)
			stmt, err := p.parseMacro()
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

	token = p.pop()
	if token.Type != lexer.ItemExpression {
		return nil, p.unexpected_token(token)
	}
	expr := &ast.Expression{
		ast.Info{Line: token.Line, Column: token.Column},
		token.Value,
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

	token = p.pop()
	if token.Type != lexer.ItemExpression {
		return nil, p.unexpected_token(token)
	}
	expr := &ast.Expression{
		ast.Info{Line: token.Line, Column: token.Column},
		token.Value,
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

func (p *parser) parseMacro() (ast.Statement, error) {
	var token lexer.Item

	token = p.pop()
	if token.Type != lexer.ItemLeftMetaColon {
		return nil, p.unexpected_token(token)
	}
	info := ast.Info{Line: token.Line, Column: token.Column}

	token = p.pop()
	if token.Type != lexer.ItemIdentifier {
		return nil, p.unexpected_token(token)
	}
	macro := &ast.Identifier{
		ast.Info{Line: token.Line, Column: token.Column},
		token.Value,
	}

	token = p.pop()
	if token.Type != lexer.ItemExpression {
		return nil, p.unexpected_token(token)
	}
	expr := &ast.Expression{
		ast.Info{Line: token.Line, Column: token.Column},
		token.Value,
	}

	if token := p.pop(); token.Type != lexer.ItemRightMeta2 {
		return nil, p.unexpected_token(token)
	}

	return &ast.Macro{info, macro, expr}, nil
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
	macro := &ast.Identifier{
		ast.Info{Line: token.Line, Column: token.Column},
		token.Value,
	}

	token = p.pop()
	if token.Type != lexer.ItemExpression {
		return nil, p.unexpected_token(token)
	}
	expr := &ast.Expression{
		ast.Info{Line: token.Line, Column: token.Column},
		token.Value,
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
	if token.Value != "end" && token.Value != macro.Value {
		return nil, fmt.Errorf("Unmatched block close tag: {{/%s}}", token.Value)
	}

	return &ast.Block{info, macro, expr, t_main, t_else}, nil
}

func (p *parser) unexpected_token(token lexer.Item) error {
	return fmt.Errorf("%d:%d unexpected token: '%s' (%v)", token.Line, token.Column, token.Value, token.Type)
}

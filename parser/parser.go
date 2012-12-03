package parser

type Context int

const (
	TextContext Context = iota
	AttrContext
	AttrTextContext
)

type Template []Statement

type Statement interface {
	Visit(Visitor)
}

type Expression interface {
}

type Visitor interface {
	EnterTemplate()
	LeaveTemplate()

	EnterBlock()
	LeaveBlock()

	VisitComment()
	VisitHelper()
	VisitValue()
	VisitLiteral()
}

type Block struct {
	expression   Expression
	template     Template
	elseTemplate Template
}

type Interpolation struct {
	expression Expression
	raw        bool
	context    Context
}

type Comment struct {
	content string
}

type Literal struct {
	content string
}

func Parse(name, content string) (Template, error) {
	p := &parser{tokenChan: lexer.Lex(name, content)}
	return p.parseTemplate(tokens)
}

type parser struct {
	tokenChan <-chan lexer.Item
	backlog   lexer.Item
}

func (p *parser) pop() lexer.Item {
	if len(p.backlog) > 0 {
		token := p.backlog[len(p.backlog)-1]
		p.backlog = p.backlog[:len(p.backlog)-1]
		return token
	}

	token, closed := <-p.tokenChan
	if closed {
		return nil
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

	for {
		token := p.pop()

		switch token.Type {

		case lexer.ItemEOF:
			return tmpl, nil

		case lexer.ItemLeftMetaSlash:
			p.push(token)
			return tmpl, nil

		case lexer.ItemHtmlText:
			tmpl = append(tmpl, Literal{token.Value})

		case lexer.ItemHtmlAttr:
			p.push(token)
			p.parseInsideHtmlAttr()

			tmpl = append(tmpl, Literal{"="})
			expr, err := p.parseExpression(tokens)
			if err != nil {
				return nil, err
			}
			tmp = append(tmpl, Interpolation{expr, false, AttrContext})

		}
	}

	return nil, fmt.Errof("Unexpected EOF")
}

package compiler

import (
	"github.com/fd/w/template/ast"
)

func (ctx *Context) CleanTemplates() {
	for _, render := range ctx.RenderFuncs {
		ast.Walk(&compact_html_literals{}, render.Template)
	}
}

type compact_html_literals struct {
}

func (visitor *compact_html_literals) Visit(n ast.Node) ast.Visitor {

	switch v := n.(type) {
	case *ast.Template:
		visitor.Perform(v)
		return visitor

	case *ast.Block:
		return visitor

	}

	return nil
}

func (visitor *compact_html_literals) Perform(tmpl *ast.Template) {
	var last_literal *ast.Literal

	new_statements := make([]ast.Statement, 0, len(tmpl.Statements))

	for _, stmt := range tmpl.Statements {
		if l, ok := stmt.(*ast.Literal); ok {
			if last_literal == nil {
				last_literal = l
				new_statements = append(new_statements, last_literal)
			} else {
				last_literal.Content += l.Content
			}

		} else {
			last_literal = nil
			new_statements = append(new_statements, stmt)

		}
	}

	tmpl.Statements = new_statements
}

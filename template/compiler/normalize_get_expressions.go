package compiler

import (
	"github.com/fd/w/template/ast"
)

func (ctx *Context) NormalizeGetExpresions() {
	for _, render := range ctx.RenderFuncs {
		ast.Walk(&normalize_get_expressions{}, render.Template)
	}
}

type normalize_get_expressions struct {
}

func (visitor *normalize_get_expressions) Visit(n ast.Node) ast.Visitor {

	switch v := n.(type) {
	case *ast.Get:
		if ident, ok := v.From.(*ast.Identifier); ok && ident.Value != "$" {
			v.From = &ast.Get{
				Info: ident.Info,
				From: &ast.Identifier{Info: ident.Info, Value: "$"},
				Name: ident,
			}
		}

	case *ast.Block:
		if ident, ok := v.Expression.(*ast.Identifier); ok && ident.Value != "$" {
			v.Expression = &ast.Get{
				Info: ident.Info,
				From: &ast.Identifier{Info: ident.Info, Value: "$"},
				Name: ident,
			}
		}

	case *ast.Interpolation:
		if ident, ok := v.Expression.(*ast.Identifier); ok && ident.Value != "$" {
			v.Expression = &ast.Get{
				Info: ident.Info,
				From: &ast.Identifier{Info: ident.Info, Value: "$"},
				Name: ident,
			}
		}

	case *ast.FunctionCall:
		if ident, ok := v.From.(*ast.Identifier); ok && ident.Value != "$" {
			v.From = &ast.Get{
				Info: ident.Info,
				From: &ast.Identifier{Info: ident.Info, Value: "$"},
				Name: ident,
			}
		}

		for i, arg := range v.Args {
			if ident, ok := arg.(*ast.Identifier); ok && ident.Value != "$" {
				v.Args[i] = &ast.Get{
					Info: ident.Info,
					From: &ast.Identifier{Info: ident.Info, Value: "$"},
					Name: ident,
				}
			}
		}

		for name, opt := range v.Options {
			if ident, ok := opt.(*ast.Identifier); ok && ident.Value != "$" {
				v.Options[name] = &ast.Get{
					Info: ident.Info,
					From: &ast.Identifier{Info: ident.Info, Value: "$"},
					Name: ident,
				}
			}
		}

	}

	return visitor
}

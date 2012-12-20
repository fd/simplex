package compiler

import (
	"github.com/fd/w/template/ast"
)

func (ctx *Context) SpecialHelpers() {
	for _, render := range ctx.RenderFuncs {
		ast.Walk(&special_helpers{ctx, render.ImportPath}, render.Template)
	}
}

type special_helpers struct {
	ctx         *Context
	import_path string
}

func (visitor *special_helpers) Visit(n ast.Node) ast.Visitor {

	switch v := n.(type) {
	case *ast.Block:
		if v.Expression != nil {
			new_expr, ok := visitor.ReplaceExpr(v.Expression)
			if ok {
				v.Expression = new_expr
			}
		}

	case *ast.Interpolation:
		if v.Expression != nil {
			new_expr, ok := visitor.ReplaceExpr(v.Expression)
			if ok {
				v.Expression = new_expr
			}
		}

	case *ast.FunctionCall:
		if v.From != nil {
			new_expr, ok := visitor.ReplaceExpr(v.From)
			if ok {
				v.From = new_expr
			}
		}

		visitor.NormalizeFunc(v)

	case *ast.Get:
		if v.From != nil {
			new_expr, ok := visitor.ReplaceExpr(v.From)
			if ok {
				v.From = new_expr
			}
		}

	}

	return visitor
}

func (visitor *special_helpers) NormalizeFunc(f *ast.FunctionCall) {
	switch f.Name {

	case "render":
		i := visitor.ctx.ImportsFor(visitor.import_path)
		pkg_name := i.Register("github.com/fd/w/runtime")
		f.Name = pkg_name + ".Render"

	case "yield":
		i := visitor.ctx.ImportsFor(visitor.import_path)
		pkg_name := i.Register("github.com/fd/w/runtime")
		f.Name = pkg_name + ".Yield"

	}
}

func (visitor *special_helpers) ReplaceExpr(e ast.Expression) (ast.Expression, bool) {
	switch v := e.(type) {
	case *ast.Get:
		if v.Name.Value == "yield" {
			func_call := &ast.FunctionCall{Info: v.Info, From: v.From, Name: "yield"}
			return func_call, true
		}

	case *ast.Identifier:
		if v.Value == "yield" {
			func_call := &ast.FunctionCall{Info: v.Info, Name: "yield"}
			return func_call, true
		}

	}

	return nil, false
}

/*

  each   (data.View,  each_branch, none_branch  template.RenderFunc)

    Posts:
    {{#each posts.All()}}
      <h1>{{ title }}</h1>
    {{else}}
      none branch
    {{/each}}

    var inner := post.All().RenderEach(render_inner_each, render_inner_none)
    var outer := w.Singleton().Render(render_outer)

    func render_outer(ctx data.Value) []Chunk {
      out := []Chunk{
        Html("Posts:\n"),
        Include{view: inner},
      }

      return out
    }

    func render_inner_each(ctx data.Value) []Chunk {
      out := []Chunk{
        Html("  <h1>"),
        escape_html(ctx.Get("title")),
        Html("</h1>\n"),
      }

      return out
    }

    func render_inner_none(ctx data.Value) []Chunk {
      out := []Chunk{
        Html("    none branch\n"),
      }

      return out
    }


  each   (data.Value, each_branch, none_branch  template.RenderFunc)

    Posts:
    {{#each comments}}
      <h1>{{ title }}</h1>
    {{else}}
      none branch
    {{/each}}

    var Outer := w.Singleton().Render(render_outer)

    func render_outer(ctx data.Value) []Chunk {
      out := []Chunk{
        Html("Posts:\n"),
        Each(ctx.Get("comments"), render_inner_each, render_inner_none)
      }

      return out
    }

    func render_inner_each(ctx data.Value) []Chunk {
      out := []Chunk{
        Html("  <h1>"),
        escape_html(ctx.Get("title")),
        Html("</h1>\n"),
      }

      return out
    }

    func render_inner_none(ctx data.Value) []Chunk {
      out := []Chunk{
        Html("    none branch\n"),
      }

      return out
    }


  if     (data.View,  any_branch,  none_branch  template.RenderFunc)
  if     (data.Value, true_branch, false_branch template.RenderFunc)
  unless (A,          true_branch, false_branch template.RenderFunc)
    ==>  if(A, false_branch, true_branch)
  with   (A,          branch template.RenderFunc)
    ==>  if(A, branch, nil)

*/

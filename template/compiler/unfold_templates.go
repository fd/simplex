package compiler

import (
	"fmt"
	"github.com/fd/w/template/ast"
	"strconv"
)

type unfold_templates struct {
	render_func *RenderFunc
	ctx         *Context
	tmpl_id     int
	view_id     int
}

func (visitor *unfold_templates) Visit(n ast.Node) ast.Visitor {
	switch v := n.(type) {
	case *ast.Template:
		visitor.perform(v)
		return visitor

	case *ast.Block:
		return visitor

	}

	return nil
}

func (visitor *unfold_templates) perform(t *ast.Template) {
	for i, stmt := range t.Statements {
		if b, ok := stmt.(*ast.Block); ok {
			incl := visitor.unfold(b)
			t.Statements[i] = incl
		}
	}
}

func (visitor *unfold_templates) unfold(b *ast.Block) ast.Statement {
	var a_branch, b_branch *RenderFunc

	a_branch = visitor.make_render_func(b.Template)

	if b.ElseTemplate != nil {
		b_branch = visitor.make_render_func(b.ElseTemplate)
	}

	view := visitor.make_data_view(b.Expression, a_branch, b_branch)
	stmt := &Include{Info: b.NodeInfo(), View: view}

	// continue unfolding a_branch and b_branch
	ast.Walk(&unfold_templates{a_branch, visitor.ctx, 0, 0}, a_branch.Template)
	if b_branch != nil {
		ast.Walk(&unfold_templates{b_branch, visitor.ctx, 0, 0}, b_branch.Template)
	}

	return stmt
}

func (visitor *unfold_templates) make_data_view(expr ast.Expression, a_branch, b_branch *RenderFunc) *DataView {
	ctx := visitor.ctx
	visitor.view_id += 1

	expr = &ast.FunctionCall{
		Info: expr.NodeInfo(),
		From: expr,
		Name: "Render___",
		Args: []ast.Expression{
			&ast.Identifier{Value: a_branch.Name()},
			&ast.Identifier{Value: b_branch.Name()},
		},
	}

	base := visitor.render_func.Name
	import_path := visitor.render_func.ImportPath
	view := &DataView{
		Name:       base + "_" + strconv.Itoa(visitor.view_id),
		ImportPath: import_path,
		Expression: expr,
	}

	name := fmt.Sprintf("\"%s\".%s#%d", import_path, base, visitor.view_id)
	ctx.DataViews[name] = view
	return view
}

func (visitor *unfold_templates) make_render_func(t *ast.Template) *RenderFunc {
	ctx := visitor.ctx
	visitor.tmpl_id += 1

	base := visitor.render_func.Name
	import_path := visitor.render_func.ImportPath
	render_func := &RenderFunc{
		Name:       base + "_" + strconv.Itoa(visitor.tmpl_id),
		ImportPath: import_path,
		Template:   t,
	}

	name := fmt.Sprintf("\"%s\".%s#%d", import_path, base, visitor.tmpl_id)
	ctx.RenderFuncs[name] = render_func
	return render_func
}

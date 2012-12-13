package compiler

import (
	"fmt"
	go_ast "go/ast"
)

func (ctx *Context) GolangFindFunctions() {
	for path, pkg := range ctx.go_ast_packages {
		v := &golang_find_functions{
			ctx:          ctx,
			package_name: path,
		}
		go_ast.Walk(v, pkg)
	}
}

type golang_find_functions struct {
	ctx          *Context
	package_name string
}

func (v *golang_find_functions) Visit(n go_ast.Node) go_ast.Visitor {
	if f, ok := n.(*go_ast.FuncDecl); ok {
		v.AnalyzeFunc(f)
		return nil
	}
	return v
}

func (v *golang_find_functions) AnalyzeFunc(f *go_ast.FuncDecl) {

	// helpers must have no receiver
	if f.Recv != nil {
		return
	}

	// helpers must be exported
	if !f.Name.IsExported() {
		return
	}

	fullname := fmt.Sprintf("\"%s\".%s", v.package_name, f.Name.String())
	v.ctx.Helpers[fullname] = f
}

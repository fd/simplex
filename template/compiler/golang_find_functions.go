package compiler

import (
	"fmt"
	go_ast "go/ast"
	go_build "go/build"
)

func (ctx *Context) GolangFindFunctions() {
	for path, pkg := range ctx.go_ast_packages {
		v := &golang_find_functions{
			ctx:          ctx,
			package_name: path,
			pkg:          pkg,
		}
		go_ast.Walk(v, pkg)
	}
}

type golang_find_functions struct {
	ctx          *Context
	package_name string
	pkg          *go_ast.Package
}

func (v *golang_find_functions) Visit(n go_ast.Node) go_ast.Visitor {
	if f, ok := n.(*go_ast.FuncDecl); ok {
		v.AnalyzeFunc(f)
		return nil
	}
	return v
}

func (v *golang_find_functions) AnalyzeFunc(f *go_ast.FuncDecl) {

	if v.IsHelperFunction(f) {
		fullname := fmt.Sprintf("\"%s\".%s", v.package_name, f.Name.String())
		v.ctx.Helpers[fullname] = f
	}

	if v.IsDataViewFunction(f) {
		fullname := fmt.Sprintf("\"%s\".%s", v.package_name, f.Name.String())
		v.ctx.DataViews[fullname] = &DataView{
			Name:       f.Name.String(),
			ImportPath: v.package_name,
		}
	}

}

func (v *golang_find_functions) IsDataViewFunction(f *go_ast.FuncDecl) bool {

	// dataviews must have no receiver
	if f.Recv != nil {
		return false
	}

	// dataviews must be exported
	if !f.Name.IsExported() {
		return false
	}

	// data views must return ([View])
	if f.Type.Results == nil {
		return false
	}

	l := f.Type.Results.List
	if len(l) != 1 {
		return false
	}

	if selector, ok := l[0].Type.(*go_ast.SelectorExpr); ok {
		pkg_name := ""
		member := selector.Sel.Name

		if ident, ok := selector.X.(*go_ast.Ident); ok && ident != nil {
			pkg_name = ident.Name

			for _, i := range v.pkg.Imports {
				if i.Name == pkg_name {
					pkg_name = i.Decl.(*go_build.Package).ImportPath
					break
				}
			}
		}

		if member != "View" || pkg_name != "github.com/fd/w/data" {
			return false
		}
	}

	return true
}

func (v *golang_find_functions) IsHelperFunction(f *go_ast.FuncDecl) bool {

	// helpers must have no receiver
	if f.Recv != nil {
		return false
	}

	// helpers must be exported
	if !f.Name.IsExported() {
		return false
	}

	// helpers must return ([Value], error) or ([Value])
	if f.Type.Results == nil {
		return false
	}

	l := f.Type.Results.List
	if len(l) < 1 || len(l) > 2 {
		return false
	}

	if selector, ok := l[0].Type.(*go_ast.SelectorExpr); ok {
		pkg_name := ""
		member := selector.Sel.Name

		if ident, ok := selector.X.(*go_ast.Ident); ok && ident != nil {
			pkg_name = ident.Name

			for _, i := range v.pkg.Imports {
				if i.Name == pkg_name {
					pkg_name = i.Decl.(*go_build.Package).ImportPath
					break
				}
			}
		}

		if member != "Value" || pkg_name != "github.com/fd/w/data" {
			return false
		}
	} else if ident, ok := l[0].Type.(*go_ast.Ident); ok {
		n := ident.Name
		if n != "string" && n != "int" && n != "float64" && n != "bool" {
			return false
		}
	} else {
		return false
	}

	if len(l) == 2 {
		if ident, ok := l[1].Type.(*go_ast.Ident); !ok {
			return false
		} else if ident.Name != "error" {
			return false
		}
	}

	return true
}

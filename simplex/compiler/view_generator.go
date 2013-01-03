package compiler

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/scanner"
)

func (pkg *Package) GenerateViews() error {
	if len(pkg.SimplexFiles) == 0 {
		return nil
	}

	generated_scope := ast.NewScope(nil)

	pkg.GeneratedFile = &ast.File{
		Name:  ast.NewIdent(pkg.AstPackage.Name),
		Scope: generated_scope,
	}
	pkg.Files["smplx_generated.go"] = pkg.GeneratedFile
	generated := &bytes.Buffer{}
	fmt.Fprintf(generated, "package %s\n", pkg.BuildPackage.Name)
	fmt.Fprintf(generated, "import sx_runtime \"github.com/fd/w/simplex/runtime\"\n")

	// resolve all simplex calls with dummy types and functions
	for _, name := range pkg.SimplexFiles {
		file := pkg.AstPackage.Files[name]
		v := &view_generator{&visitor{file, pkg, nil}, generated, map[string]bool{}}
		ast.Walk(v, file)
		err := v.errors.Err()
		if err != nil {
			return err
		}
	}

	f, err := parser.ParseFile(
		pkg.FileSet,
		"smplx_generated.go",
		generated,
		parser.SpuriousErrors|parser.ParseComments,
	)
	if err != nil {
		return err
	}

	for _, obj := range generated_scope.Objects {
		new_obj := f.Scope.Lookup(obj.Name)
		if new_obj != nil {
			new_obj.Data = obj.Data
		}
	}

	pkg.Files["smplx_generated.go"] = f

	return nil
}

func (pkg *Package) ResolveViews() error {
	// resolve all simplex calls with dummy types and functions
	for _, name := range pkg.SimplexFiles {
		file := pkg.AstPackage.Files[name]
		v := &view_resolver{&visitor{file, pkg, nil}}
		ast.Walk(v, file)
		err := v.errors.Err()
		if err != nil {
			return err
		}
	}

	return nil
}

/*
  Declare new View types
  source(Type)         => TypeView
  collect(func(T)Type) => TypeView
*/

type view_generator struct {
	*visitor
	Generated *bytes.Buffer
	Views     map[string]bool
}

func (v *view_generator) Write(dat []byte) (int, error) {
	return v.Generated.Write(dat)
}

func (v *view_generator) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.CallExpr:
		return v.VisitCallExpr(n)
	}
	return v
}

func (v *view_generator) VisitCallExpr(call *ast.CallExpr) ast.Visitor {
	switch node := call.Fun.(type) {

	case *ast.Ident:
		switch node.Name {
		case "source": // source(...)
			return v.VisitSourceCall(call)
		}

	case *ast.SelectorExpr:
		switch node.Sel.Name {
		case "collect": // x.collect(...)
			return v.VisitCollectCall(call)
		}

	}

	return v
}

func (v *view_generator) VisitSourceCall(call *ast.CallExpr) ast.Visitor {
	if len(call.Args) != 1 {
		v.push_errorf(call, "Expected a local type")
		return nil
	}

	typ := call.Args[0]
	ident, _ := typ.(*ast.Ident)

	if ident == nil {
		v.push_errorf(call.Args[0], "Expected a local type")
		return nil
	}

	if ident.Obj == nil {
		v.push_errorf(call.Args[0], "Expected a local type")
		return nil
	}

	view_obj, err := v.Pkg.declareView(ident.Name)
	if err != nil {
		v.push_error(ident, err)
		return nil
	}

	call.Fun.(*ast.Ident).Obj = view_obj
	v.PrintTypeDecl(view_obj)

	return v
}

func (v *view_generator) VisitCollectCall(call *ast.CallExpr) ast.Visitor {
	if len(call.Args) != 1 {
		v.push_errorf(call, "Expected function with 1 return value")
		return nil
	}

	collect_func, err := find_inner_function(call.Args[0], v.File.Scope)
	if err != nil {
		v.push_error(call.Args[0], err)
		return nil
	}

	if collect_func.Type.Results.NumFields() != 1 {
		v.push_errorf(call.Args[0], "Expected function with 1 return value")
		return nil
	}

	typ := collect_func.Type.Results.List[0].Type
	ident, ok := typ.(*ast.Ident)
	if !ok {
		v.push_errorf(call.Args[0], "Expected function with a local return value")
		return nil
	}
	if ident.Obj == nil {
		v.push_errorf(call.Args[0], "Expected function with a local return value")
		return nil
	}

	view_obj, err := v.Pkg.declareView(ident.Obj.Name)
	if err != nil {
		v.push_error(ident, err)
		return nil
	}

	call.Fun.(*ast.SelectorExpr).Sel.Obj = view_obj
	v.PrintTypeDecl(view_obj)

	return v
}

func (v *view_generator) PrintTypeDecl(view_obj *ast.Object) {
	if !v.Views[view_obj.Name] {
		v.Views[view_obj.Name] = true
		fmt.Fprintf(v, `type %s struct { view sx_runtime.View }
`, view_obj.Name)
		fmt.Fprintf(v, `func %sSource()%s{return %s{ sx_runtime.Source("%s") }}
`,
			view_obj.Name,
			view_obj.Name,
			view_obj.Name,
			view_obj.Name)
		fmt.Fprintf(v, `func (w %s)Where(f sx_runtime.WhereFunc)%s{return %s{ w.view.Where(f) }}
`,
			view_obj.Name,
			view_obj.Name,
			view_obj.Name)
		fmt.Fprintf(v, `func (w %s)Sort(f sx_runtime.SortFunc)%s{return %s{ w.view.Sort(f) }}
`,
			view_obj.Name,
			view_obj.Name,
			view_obj.Name)
		fmt.Fprintf(v, `func (w %s)Group(f sx_runtime.GroupFunc)%s{return %s{ w.view.Group(f) }}
`,
			view_obj.Name,
			view_obj.Name,
			view_obj.Name)
		fmt.Fprintf(v, `func %sCollectedFrom(input sx_runtime.ViewWrapper, f sx_runtime.CollectFunc)%s{return %s{ input.View().Collect(f) }}
`,
			view_obj.Name,
			view_obj.Name,
			view_obj.Name)
		fmt.Fprintf(v, `func (w %s)View()sx_runtime.View{return w.view }
`,
			view_obj.Name)
	}
}

type view_resolver struct {
	*visitor
}

func (v *view_resolver) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.CallExpr:
		return v.VisitCallExpr(n)
	}
	return v
}

func (v *view_resolver) VisitCallExpr(call *ast.CallExpr) ast.Visitor {
	switch node := call.Fun.(type) {

	case *ast.Ident:
		switch node.Name {
		case "source": // source(...)
			return v.VisitSourceCall(call)
		}

	case *ast.SelectorExpr:
		switch node.Sel.Name {

		case "where": // x.where(...)
			return v.VisitWhereCall(call)

		case "collect": // x.collect(...)
			return v.VisitCollectCall(call)

		case "sort": // x.where(...)
			return v.VisitSortCall(call)

		case "group": // x.where(...)
			return v.VisitGroupCall(call)

		}

	}

	return v
}

func (v *view_resolver) VisitSourceCall(call *ast.CallExpr) ast.Visitor {
	view_obj := resolve_type(v.File.Scope, call)
	view_typ := view_obj.Type.(*ViewDecl).MemberType

	fident := ast.NewIdent(view_typ + "ViewSource")
	fident.Obj = view_obj
	call.Args = []ast.Expr{}
	call.Fun = fident

	return v
}

func (v *view_resolver) VisitWhereCall(call *ast.CallExpr) ast.Visitor {
	ast.Walk(v, call.Fun)

	view_obj := resolve_type(v.File.Scope, call)
	view_typ := view_obj.Type.(*ViewDecl).MemberType
	inner := call.Args[0]

	wrapper := &ast.FuncLit{
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					&ast.Field{
						Names: []*ast.Ident{
							ast.NewIdent("m"),
						},
						Type: &ast.InterfaceType{
							Methods: &ast.FieldList{Opening: call.Pos(), Closing: call.Pos()},
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					&ast.Field{
						Type: ast.NewIdent("bool"),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.CallExpr{
							Fun: inner,
							Args: []ast.Expr{
								&ast.TypeAssertExpr{
									X:    ast.NewIdent("m"),
									Type: ast.NewIdent(view_typ),
								},
							},
						},
					},
				},
			},
		},
	}

	call.Args = []ast.Expr{wrapper}
	fident := ast.NewIdent("Where")
	fident.Obj = view_obj
	call.Fun.(*ast.SelectorExpr).Sel = fident

	return nil
}

func (v *view_resolver) VisitCollectCall(call *ast.CallExpr) ast.Visitor {
	ast.Walk(v, call.Fun)

	i_obj := resolve_type(v.File.Scope, call.Fun.(*ast.SelectorExpr).X)
	o_obj := resolve_type(v.File.Scope, call)

	i_typ := i_obj.Type.(*ViewDecl).MemberType
	o_typ := o_obj.Type.(*ViewDecl).MemberType + "View"

	first := call.Fun.(*ast.SelectorExpr).X
	inner := call.Args[0]

	wrapper := &ast.FuncLit{
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					&ast.Field{
						Names: []*ast.Ident{
							ast.NewIdent("m"),
						},
						Type: &ast.InterfaceType{
							Methods: &ast.FieldList{Opening: call.Pos(), Closing: call.Pos()},
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					&ast.Field{
						Type: &ast.InterfaceType{
							Methods: &ast.FieldList{Opening: call.Pos(), Closing: call.Pos()},
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.CallExpr{
							Fun: inner,
							Args: []ast.Expr{
								&ast.TypeAssertExpr{
									X:    ast.NewIdent("m"),
									Type: ast.NewIdent(i_typ),
								},
							},
						},
					},
				},
			},
		},
	}

	call.Args = []ast.Expr{first, wrapper}
	fident := ast.NewIdent(o_typ + "CollectedFrom")
	fident.Obj = o_obj
	call.Fun = fident

	return nil
}

func (v *view_resolver) VisitSortCall(call *ast.CallExpr) ast.Visitor {
	ast.Walk(v, call.Fun)

	view_obj := resolve_type(v.File.Scope, call)
	view_typ := view_obj.Type.(*ViewDecl).MemberType
	inner := call.Args[0]

	wrapper := &ast.FuncLit{
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					&ast.Field{
						Names: []*ast.Ident{
							ast.NewIdent("m"),
						},
						Type: &ast.InterfaceType{
							Methods: &ast.FieldList{Opening: inner.Pos(), Closing: inner.Pos()},
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					&ast.Field{
						Type: &ast.InterfaceType{
							Methods: &ast.FieldList{Opening: inner.Pos(), Closing: inner.Pos()},
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.CallExpr{
							Fun: inner,
							Args: []ast.Expr{
								&ast.TypeAssertExpr{
									X:    ast.NewIdent("m"),
									Type: ast.NewIdent(view_typ),
								},
							},
						},
					},
				},
			},
		},
	}

	call.Args = []ast.Expr{wrapper}
	fident := ast.NewIdent("Sort")
	fident.Obj = view_obj
	call.Fun.(*ast.SelectorExpr).Sel = fident

	return nil
}

func (v *view_resolver) VisitGroupCall(call *ast.CallExpr) ast.Visitor {
	ast.Walk(v, call.Fun)

	view_obj := resolve_type(v.File.Scope, call)
	view_typ := view_obj.Type.(*ViewDecl).MemberType
	inner := call.Args[0]

	wrapper := &ast.FuncLit{
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					&ast.Field{
						Names: []*ast.Ident{
							ast.NewIdent("m"),
						},
						Type: &ast.InterfaceType{
							Methods: &ast.FieldList{Opening: call.Pos(), Closing: call.Pos()},
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					&ast.Field{
						Type: &ast.InterfaceType{
							Methods: &ast.FieldList{Opening: call.Pos(), Closing: call.Pos()},
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.CallExpr{
							Fun: inner,
							Args: []ast.Expr{
								&ast.TypeAssertExpr{
									X:    ast.NewIdent("m"),
									Type: ast.NewIdent(view_typ),
								},
							},
						},
					},
				},
			},
		},
	}

	call.Args = []ast.Expr{wrapper}
	fident := ast.NewIdent("Group")
	fident.Obj = view_obj
	call.Fun.(*ast.SelectorExpr).Sel = fident

	return nil
}

/*
  find a literal function reference
*/
func find_inner_function(node ast.Node, scope *ast.Scope) (*ast.FuncDecl, error) {
	var obj *ast.Object

	switch n := node.(type) {
	case *ast.Ident:
		if n.Obj == nil {
			resolve(scope, n)
		}
		obj = n.Obj

	case *ast.SelectorExpr:
		if n.Sel.Obj == nil {
			if ident, ok := n.X.(*ast.Ident); ok && ident.Obj.Type == ast.Pkg {
				resolve(ident.Obj.Data.(*ast.Scope), n.Sel)
			}
		}
		obj = n.Sel.Obj

	default:
		return nil, fmt.Errorf("Expected a function reference.")

	}

	if obj == nil {
		return nil, fmt.Errorf("Expected a function reference.")
	}
	if obj.Kind != ast.Fun {
		return nil, fmt.Errorf("Expected a function reference.")
	}

	func_decl, _ := obj.Decl.(*ast.FuncDecl)
	return func_decl, nil
}

type visitor struct {
	File   *ast.File
	Pkg    *Package
	errors scanner.ErrorList
}

func (v *visitor) push_error_msg(node ast.Node, msg string) {
	pos := v.Pkg.FileSet.Position(node.Pos())
	v.errors.Add(pos, msg)
}

func (v *visitor) push_error(node ast.Node, msg error) {
	v.push_error_msg(node, msg.Error())
}

func (v *visitor) push_errorf(node ast.Node, msg string, arg ...interface{}) {
	v.push_error_msg(node, fmt.Sprintf(msg, arg...))
}

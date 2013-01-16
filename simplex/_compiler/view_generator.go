package compiler

import (
	"bytes"
	"fmt"
	sx_ast "github.com/fd/w/simplex/ast"
	"go/ast"
	"go/parser"
	"go/scanner"
	"go/token"
)

func (pkg *Package) GenerateViews() error {
	if len(pkg.SmplxFiles) == 0 {
		return nil
	}

	generated_scope := ast.NewScope(nil)

	pkg.Views = map[string]*ViewDecl{}
	pkg.GeneratedFile = &ast.File{
		Name:  ast.NewIdent(pkg.AstPackage.Name),
		Scope: generated_scope,
	}
	pkg.Files["smplx_generated.go"] = pkg.GeneratedFile
	generated := &bytes.Buffer{}
	fmt.Fprintf(generated, "package %s\n", pkg.BuildPackage.Name)
	fmt.Fprintf(generated, "import sx_runtime \"github.com/fd/w/simplex/runtime\"\n")

	// resolve all simplex calls with dummy types and functions
	for _, file := range pkg.SmplxFiles {
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

	// Link view functions
	for _, e := range f.Decls {
		switch decl := e.(type) {
		case *ast.FuncDecl:
			if decl.Recv == nil {
				continue
			}

			recv_type := decl.Recv.List[0].Type
			ident, ok := recv_type.(*ast.Ident)
			if !ok {
				continue
			}

			view_decl, ok := pkg.Views[ident.Name]
			if !ok {
				continue
			}

			switch decl.Name.Name {
			case "Source":
				view_decl.Source = decl
			case "Select":
				view_decl.Select = decl
			case "Reject":
				view_decl.Reject = decl
			case "Sort":
				view_decl.Sort = decl
			case "Group":
				view_decl.Group = decl
			case "CollectedFrom":
				view_decl.CollectedFrom = decl
			default:
				continue
			}

			decl.Name.Obj = ast.NewObj(ast.Fun, ident.Name+"•"+decl.Name.Name)
			decl.Name.Obj.Decl = decl
			f.Scope.Insert(decl.Name.Obj)

		case *ast.GenDecl:
			if decl.Tok != token.TYPE {
				continue
			}

			if len(decl.Specs) != 1 {
				continue
			}

			type_spec, ok := decl.Specs[0].(*ast.TypeSpec)
			if !ok {
				continue
			}

			if type_spec.Name == nil {
				continue
			}

			view_decl, ok := pkg.Views[type_spec.Name.Name]
			if !ok {
				continue
			}

			view_decl.ViewType = type_spec.Name
			type_spec.Name.Obj.Data = view_decl
		}
	}

	pkg.Files["smplx_generated.go"] = f

	return nil
}

func (pkg *Package) ResolveViews() error {
	// resolve all simplex calls with dummy types and functions
	for _, file := range pkg.SmplxFiles {
		v := &view_resolver{&visitor{file, pkg, nil}}
		ast.Walk(v, file)
		err := v.errors.Err()
		if err != nil {
			return err
		}
	}

	return nil
}

type view_generator struct {
	*visitor
	Generated *bytes.Buffer
	Views     map[string]bool
}

func (v *view_generator) Write(dat []byte) (int, error) {
	return v.Generated.Write(dat)
}

func (v *view_generator) Visit(node ast.Node) ast.Visitor {
	decl, ok := node.(*ast.GenDecl)
	if !ok {
		return v
	}

	if decl.Tok != token.TYPE {
		return v
	}

	for _, any_spec := range decl.Specs {
		spec := any_spec.(*ast.TypeSpec)

		if spec.Name == nil {
			continue
		}

		struc, ok := spec.Type.(*ast.StructType)
		if !ok {
			continue
		}

		if len(struc.Fields.List) < 1 {
			continue
		}

		field := struc.Fields.List[0]
		if len(field.Names) != 0 {
			continue
		}

		ident, ok := field.Type.(*ast.Ident)
		if !ok {
			continue
		}

		if ident.Name != "view" {
			continue
		}

		struc.Fields.List = struc.Fields.List[1:]

		view_decl := &ViewDecl{
			MemberType: spec.Name,
		}
		v.Pkg.Views[spec.Name.Name+"View"] = view_decl
		spec.Name.Obj.Data = view_decl

		v.PrintTypeDecl(spec.Name)
	}

	return nil
}

func (v *view_generator) PrintTypeDecl(view_ident *ast.Ident) {
	decl := view_ident.Obj.Data.(*ViewDecl)
	m_name := decl.MemberType.Name
	v_name := m_name + "View"

	if !v.Views[v_name] {
		v.Views[v_name] = true

		fmt.Fprintf(v, `type %s struct { view sx_runtime.View }
`, v_name)
		fmt.Fprintf(v, `func (w %s)Source()%s{return %s{ sx_runtime.Source("%s") }}
`,
			v_name,
			v_name,
			v_name,
			v_name)
		fmt.Fprintf(v, `func (w %s)Select(f sx_runtime.SelectFunc)%s{return %s{ w.view.Select(f) }}
`,
			v_name,
			v_name,
			v_name)
		fmt.Fprintf(v, `func (w %s)Sort(f sx_runtime.SortFunc)%s{return %s{ w.view.Sort(f) }}
`,
			v_name,
			v_name,
			v_name)
		fmt.Fprintf(v, `func (w %s)Group(f sx_runtime.GroupFunc)%s{return %s{ w.view.Group(f) }}
`,
			v_name,
			v_name,
			v_name)
		fmt.Fprintf(v, `func (w %s)CollectedFrom(input sx_runtime.ViewWrapper, f sx_runtime.CollectFunc)%s{return %s{ input.View().Collect(f) }}
`,
			v_name,
			v_name,
			v_name)
		fmt.Fprintf(v, `func (w %s)View()sx_runtime.View{return w.view }
`,
			v_name)
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

		case "select": // x.select(...)
			return v.VisitSelectCall(call)

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
	if len(call.Args) != 1 {
		return v
	}

	ident, ok := call.Args[0].(*ast.Ident)
	if !ok {
		return v
	}

	view_decl, ok := v.Pkg.Views[ident.Name+"View"]
	if !ok {
		return v
	}

	call.Args = []ast.Expr{}
	call.Fun = &ast.SelectorExpr{
		X: &ast.CompositeLit{
			Type: view_decl.ViewType,
		},
		Sel: view_decl.Source.Name,
	}

	return v
}

// T.where(F) => T
func (v *view_resolver) VisitSelectCall(call *ast.CallExpr) ast.Visitor {
	ast.Walk(v, call.Fun)

	view_decl := resolve_view_type(v.File.Scope, call.Fun.(*ast.SelectorExpr).X)
	inner := call.Args[0]

	wrapper := &ast.FuncLit{
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
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
					{
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
									Type: view_decl.MemberType,
								},
							},
						},
					},
				},
			},
		},
	}

	call.Args = []ast.Expr{wrapper}
	call.Fun.(*ast.SelectorExpr).Sel = view_decl.Select.Name

	return nil
}

func (v *view_resolver) VisitCollectCall(call *ast.CallExpr) ast.Visitor {
	ast.Walk(v, call.Fun)

	first := call.Fun.(*ast.SelectorExpr).X
	inner := call.Args[0]

	//ast.Print(v.Pkg.FileSet, inner)

	i_decl := resolve_view_type(v.File.Scope, first)
	o_decl := resolve_function_type(v.File.Scope, inner)

	wrapper := &ast.FuncLit{
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
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
					{
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
									Type: i_decl.MemberType,
								},
							},
						},
					},
				},
			},
		},
	}

	call.Args = []ast.Expr{first, wrapper}
	call.Fun = &ast.SelectorExpr{
		X: &ast.CompositeLit{
			Type: o_decl.ViewType,
		},
		Sel: o_decl.CollectedFrom.Name,
	}

	return nil
}

func (v *view_resolver) VisitSortCall(call *ast.CallExpr) ast.Visitor {
	ast.Walk(v, call.Fun)

	view_decl := resolve_view_type(v.File.Scope, call.Fun.(*ast.SelectorExpr).X)
	inner := call.Args[0]

	wrapper := &ast.FuncLit{
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
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
					{
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
									Type: view_decl.MemberType,
								},
							},
						},
					},
				},
			},
		},
	}

	call.Args = []ast.Expr{wrapper}
	call.Fun.(*ast.SelectorExpr).Sel = view_decl.Sort.Name

	return nil
}

func (v *view_resolver) VisitGroupCall(call *ast.CallExpr) ast.Visitor {
	ast.Walk(v, call.Fun)

	view_decl := resolve_view_type(v.File.Scope, call.Fun.(*ast.SelectorExpr).X)
	inner := call.Args[0]

	wrapper := &ast.FuncLit{
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
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
					{
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
									Type: view_decl.MemberType,
								},
							},
						},
					},
				},
			},
		},
	}

	call.Args = []ast.Expr{wrapper}
	call.Fun.(*ast.SelectorExpr).Sel = view_decl.Group.Name

	return nil
}

type visitor struct {
	File   *sx_ast.File
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

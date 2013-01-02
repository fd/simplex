package compiler

import (
	"fmt"
	"go/ast"
	"go/scanner"
)

func (pkg *Package) GenerateViews() error {
	pkg.GeneratedFile = &ast.File{
		Name:  ast.NewIdent(pkg.AstPackage.Name),
		Scope: ast.NewScope(nil),
	}
	pkg.Files["smplx_generated.go"] = pkg.GeneratedFile

	// resolve all simplex calls with dummy types and functions
	for _, name := range pkg.SimplexFiles {
		file := pkg.AstPackage.Files[name]
		v := &view_generator{&visitor{file, pkg, nil}}
		ast.Walk(v, file)
		err := v.errors.Err()
		if err != nil {
			return err
		}
	}

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

	return v
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

		}

	}

	return v
}

func (v *view_resolver) VisitSourceCall(call *ast.CallExpr) ast.Visitor {
	view_obj := resolve_type(v.File.Scope, call)
	fmt.Print("source: ")
	ast.Print(v.Pkg.FileSet, view_obj.Type)
	return v
}

func (v *view_resolver) VisitWhereCall(call *ast.CallExpr) ast.Visitor {
	ast.Walk(v, call.Fun)

	view_obj := resolve_type(v.File.Scope, call)

	fident := ast.NewIdent("Where")
	fident.Obj = view_obj
	call.Fun.(*ast.SelectorExpr).Sel = fident

	return nil
}

func (v *view_resolver) VisitCollectCall(call *ast.CallExpr) ast.Visitor {
	ast.Walk(v, call.Fun)

	//i_obj := resolve_type(v.File.Scope, call.Fun.(*ast.SelectorExpr).X)
	o_obj := resolve_type(v.File.Scope, call)

	//i_typ := i_obj.Type.(*ViewDecl).MemberType + "View"
	o_typ := o_obj.Type.(*ViewDecl).MemberType + "View"

	first := call.Fun.(*ast.SelectorExpr).X
	//inner := call.Args[0]

	wrapper := &ast.FuncLit{
		Type: &ast.FuncType{},
		Body: &ast.BlockStmt{},
	}

	call.Args = []ast.Expr{first, wrapper}
	fident := ast.NewIdent(o_typ + "CollectedFrom")
	fident.Obj = o_obj
	call.Fun = fident

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

/*
func (v *view_generator) VisitWhereCall(call *ast.CallExpr) ast.Visitor {
  recv, err := find_receiver_view(v.File.Scope, call.Fun.(*ast.SelectorExpr).X)
  if err != nil {
    v.push_error(call, err)
    return nil
  }

  if len(call.Args) != 1 {
    v.push_errorf(call, "Expected a function of type func(T)bool")
    return nil
  }

  func_decl, err := find_inner_function(call.Args[0], v.Pkg.AstPackage.Scope)
  if err != nil {
    v.push_error(call.Args[0], err)
    return nil
  }

  if !verify_func_with_one_return_of_type(func_decl, "bool") {
    v.push_errorf(func_decl, "Expected a function of type func(T)bool")
    return nil
  }

  fobj, err := v.Pkg.declareWhere(recv.ViewDecl())
  if err != nil {
    v.push_error(call, err)
    return v
  }

  fident := ast.NewIdent(fobj.Name)
  fident.Obj = fobj

  call.Fun = fident
  return nil
}
*/

/*
Ident
  X.where(...)
SelectorExpr
  P.X.where(...)
CallExpr
  F(...).where(...)

X must resolve to a View as local or global variable/argument
P must be a package
F must return a single value which must be a view
*/
/*
func find_receiver_view(scope *ast.Scope, recv ast.Expr) (View, error) {
  switch n := recv.(type) {

  case *ast.Ident:
    if n.Obj == nil {
      resolve(scope, n)
    }
    obj := resolve_type(scope, n)
    ast.Print(nil, n)
    ast.Print(nil, obj)
  }

  return nil, nil
}

func verify_func_with_one_return(f *ast.FuncDecl) ast.Expr {
  l := f.Type.Results.List
  if len(l) != 1 {
    return nil
  }

  r := l[0]
  if len(r.Names) > 1 {
    return nil
  }

  return r.Type
}

func verify_func_with_one_return_of_type(f *ast.FuncDecl, typ_name string) bool {
  typ := verify_func_with_one_return(f)
  if typ == nil {
    return false
  }

  ident, _ := typ.(*ast.Ident)
  if ident == nil {
    return false
  }

  return ident.Name == typ_name
}

const ViewCode = `
type {{Type}}View struct {
  view simplex_runtime.View
}

func {{Type}}ViewSource () {{Type}}View {
  v := new({{Type}}View)
  v.view.Name = "{{Type}}"
  v.view.Type = &{{Type}}{}
  return v
}

func (v {{Type}}View) Where (f func({{Type}})bool) {{Type}}View {
  v.view = v.view.Where(func(member interface{})bool{
    if typed, ok := member.({{Type}}); ok {
      return f(typed)
    }
    return false
  })
  return v
}
`

const ViewSortCode = `
func (v {{Type}}View) Sort_{{Type2}} (f func({{Type}}){{Type2}}) {{Type}}View {
  v.view = v.view.Sort(func(member interface{})interface{}{
    if typed, ok := member.({{Type}}); ok {
      return f(typed)
    }
    return nil
  })
  return v
}
`

const ViewMapCode = `
func (v {{Type}}View) Map_{{Type2}} (f func({{Type}}){{Type2}}) {{Type}}View {
  v.view = v.view.Map(func(member interface{})interface{}{
    if typed, ok := member.({{Type}}); ok {
      return f(typed)
    }
    return nil
  })
  return v
}
`

const ViewGroupCode = `
func (v {{Type}}View) Group_{{Type2}} (f func({{Type}}){{Type2}}) {{Type}}View {
  v.view = v.view.Group(func(member interface{})interface{}{
    if typed, ok := member.({{Type}}); ok {
      return f(typed)
    }
    return nil
  })
  return v
}
`*/

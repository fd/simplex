package compiler

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/token"
)

var Context build.Context
var Packages = map[string]*Package{}

func init() {
	Context = build.Default
	Context.CgoEnabled = false
}

type Package struct {
	SimplexFiles []string
	BuildPackage *build.Package
	FileSet      *token.FileSet
	Files        map[string]*ast.File
	AstPackage   *ast.Package
}

type View interface {
	ViewDecl() *ViewDecl
}

type ViewDecl struct {
	MemberType string
	Where      *WhereDecl
}

type SourceDecl struct {
	View *ViewDecl
}

type WhereDecl struct {
	View *ViewDecl
}

func (pkg *Package) declareView(type_name string) (*ast.Object, error) {
	obj_name := type_name + "View"

	obj := pkg.AstPackage.Scope.Lookup(obj_name)
	if obj != nil {
		if _, ok := obj.Decl.(*ViewDecl); ok {
			return obj, nil
		} else {
			return nil, fmt.Errorf("%s is not a simplex.View", obj_name)
		}
	}

	obj = ast.NewObj(ast.Con, obj_name)
	obj.Decl = &ViewDecl{MemberType: type_name}
	pkg.AstPackage.Scope.Insert(obj)

	return obj, nil
}

func (pkg *Package) declareSource(type_name string) (*ast.Object, error) {
	view, err := pkg.declareView(type_name)
	if err != nil {
		return nil, err
	}

	obj_name := type_name + "ViewSource"

	obj := pkg.AstPackage.Scope.Lookup(obj_name)
	if obj != nil {
		if _, ok := obj.Decl.(*SourceDecl); ok {
			return obj, nil
		} else {
			return nil, fmt.Errorf("%s is not a simplex.ViewSource", obj_name)
		}
	}

	obj = ast.NewObj(ast.Fun, obj_name)
	obj.Decl = &SourceDecl{view.Decl.(*ViewDecl)}
	pkg.AstPackage.Scope.Insert(obj)

	return obj, nil
}

func (pkg *Package) declareWhere(view *ViewDecl) (*ast.Object, error) {
	obj_name := view.MemberType + "ViewSource"

	if view.Where == nil {
		view.Where = &WhereDecl{view}
	}

	obj := ast.NewObj(ast.Fun, "Where")
	obj.Decl = view.Where
	return obj, nil
}

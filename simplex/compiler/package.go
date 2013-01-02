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

	GeneratedFile *ast.File
}

type ViewDecl struct {
	MemberType string
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

	obj = ast.NewObj(ast.Typ, obj_name)
	obj.Data = &ViewDecl{MemberType: type_name}
	pkg.GeneratedFile.Scope.Insert(obj)

	fmt.Println("Generated view:", obj_name)

	return obj, nil
}

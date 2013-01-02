package compiler

import (
	"go/ast"
)

func (pkg *Package) ResolvePackage() error {
	if pkg.AstPackage != nil {
		return nil
	}

	importer := func(imports map[string]*ast.Object, dir string) (*ast.Object, error) {

		if dir == "unsafe" {
			return Unsafe, nil
		}

		if obj, p := imports[dir]; p {
			return obj, nil
		}

		pkg, err := ImportResolved(dir, pkg.BuildPackage.ImportPath)
		if err != nil {
			return nil, err
		}

		obj := ast.NewObj(ast.Pkg, pkg.BuildPackage.Name)
		obj.Decl = pkg.BuildPackage
		obj.Data = pkg.AstPackage.Scope
		imports[pkg.BuildPackage.ImportPath] = obj
		imports[dir] = obj

		return obj, nil
	}

	ast_pkg, _ := ast.NewPackage(pkg.FileSet, pkg.Files, importer, GoUniverse)
	pkg.AstPackage = ast_pkg

	err := pkg.GenerateViews()
	if err != nil {
		return err
	}

	ast_pkg, err = ast.NewPackage(pkg.FileSet, pkg.Files, importer, GoUniverse)
	if err != nil {
		return err
	}
	pkg.AstPackage = ast_pkg

	err = pkg.ResolveViews()
	if err != nil {
		return err
	}

	return nil
}

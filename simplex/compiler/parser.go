package compiler

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path"
	"time"
)

func (pkg *Package) ParseFiles() error {
	if pkg.Files != nil {
		return nil
	}

	pkg.FileSet = token.NewFileSet()
	pkg.Files = map[string]*ast.File{}
	pkg.ModTimes = map[string]time.Time{}

	// Go files
	for _, name := range pkg.BuildPackage.GoFiles {
		if name == "smplx_generated.go" {
			continue
		}

		l_name := path.Join(pkg.BuildPackage.ImportPath, name)
		r_name := path.Join(pkg.BuildPackage.Dir, name)

		source, err := ioutil.ReadFile(r_name)
		if err != nil {
			return err
		}

		stat, err := os.Stat(r_name)
		if err != nil {
			return err
		}

		f, err := parser.ParseFile(
			pkg.FileSet,
			l_name,
			source,
			parser.SpuriousErrors|parser.ParseComments,
		)
		if err != nil {
			return err
		}

		pkg.ModTimes[name] = stat.ModTime()
		pkg.Files[name] = f
	}

	// Simplex files
	for _, name := range pkg.SimplexFiles {
		l_name := path.Join(pkg.BuildPackage.ImportPath, name)
		r_name := path.Join(pkg.BuildPackage.Dir, name)

		source, err := ioutil.ReadFile(r_name)
		if err != nil {
			return err
		}

		stat, err := os.Stat(r_name)
		if err != nil {
			return err
		}

		f, err := parser.ParseFile(
			pkg.FileSet,
			l_name,
			source,
			parser.SpuriousErrors|parser.ParseComments,
		)
		if err != nil {
			return err
		}

		resolve_simplex_file(f)

		pkg.ModTimes[name] = stat.ModTime()
		pkg.Files[name] = f
	}

	return nil
}

func resolve_simplex_file(f *ast.File) {
	orig_outer := f.Scope.Outer
	defer func() {
		f.Scope.Outer = orig_outer
	}()

	f.Scope.Outer = SmplxUniverse

	i := 0
	for _, ident := range f.Unresolved {
		if !resolve(f.Scope, ident) {
			f.Unresolved[i] = ident
			i++
		}
	}
	f.Unresolved = f.Unresolved[0:i]
}

func resolve(scope *ast.Scope, ident *ast.Ident) bool {
	for ; scope != nil; scope = scope.Outer {
		if obj := scope.Lookup(ident.Name); obj != nil {
			ident.Obj = obj
			return true
		}
	}
	return false
}

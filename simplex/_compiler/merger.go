package compiler

import (
	"go/ast"
	"go/token"
	"time"
)

func (pkg *Package) MergeGeneratedFiles() error {
	orig_files := pkg.AstPackage.Files
	gen_files := map[string]*ast.File{}

	var (
		max_mod_time time.Time
	)

	if f, ok := orig_files["smplx_generated.go"]; ok {
		gen_files["smplx_generated.go"] = f
		delete(orig_files, "smplx_generated.go")
	}
	for name := range pkg.SmplxFiles {
		gen_files[name] = orig_files[name]
		mod_time := pkg.ModTimes[name]

		if mod_time.After(max_mod_time) {
			max_mod_time = mod_time
		}

		delete(orig_files, name)
		delete(pkg.ModTimes, name)
	}
	pkg.AstPackage.Files = gen_files

	f := ast.MergePackageFiles(pkg.AstPackage, ast.FilterImportDuplicates)
	collect_imports_at_the_top(f)
	orig_files["smplx_generated.go"] = f

	pkg.GeneratedFile = f
	pkg.AstPackage.Files = orig_files
	pkg.SmplxFiles = nil
	pkg.Files = orig_files
	pkg.ModTimes["smplx_generated.go"] = max_mod_time

	return nil
}

func collect_imports_at_the_top(f *ast.File) {
	decl := f.Decls
	imports := []ast.Decl{}

	for i := len(decl) - 1; i >= 0; i-- {
		n := decl[i]
		d, ok := n.(*ast.GenDecl)
		if !ok || d.Tok != token.IMPORT {
			continue
		}

		imports = append(imports, d)

		if i > 0 {
			decl = append(decl[:i], decl[i+1:]...)
		} else {
			decl = decl[i+1:]
		}
	}

	f.Decls = append(imports, decl...)
}

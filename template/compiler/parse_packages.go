package compiler

import (
	"fmt"
	go_ast "go/ast"
	go_build "go/build"
	go_parser "go/parser"
	"os"
)

func go_file_filter(pkg *go_build.Package) func(os.FileInfo) bool {
	return func(f os.FileInfo) bool {
		base := f.Name()

		for _, name := range pkg.GoFiles {
			if name == base {
				return true
			}
		}

		return false
	}
}

func (ctx *Context) go_package_importer(imports map[string]*go_ast.Object, dir, srcdir string) (*go_ast.Object, error) {
	var ast_pkg_obj *go_ast.Object

	defer func() {
		if ast_pkg_obj == nil {
			fmt.Printf("[%s] path: %s (not found)\n", srcdir, dir)
		}
	}()

	if dir == "unsafe" {
		ast_pkg_obj = Unsafe
		return ast_pkg_obj, nil
	}

	ast_pkg_obj = imports[dir]
	if ast_pkg_obj != nil {
		return ast_pkg_obj, nil
	}

	build_pkg, ast_pkg, err := ctx.go_import_package(dir, srcdir)
	if err != nil {
		return nil, err
	}

	ast_pkg_obj = go_ast.NewObj(go_ast.Pkg, build_pkg.Name)
	ast_pkg_obj.Decl = build_pkg
	ast_pkg_obj.Data = ast_pkg.Scope
	imports[build_pkg.ImportPath] = ast_pkg_obj
	imports[dir] = ast_pkg_obj

	return ast_pkg_obj, nil
}

func (ctx *Context) go_import_package(dir, srcdir string) (*go_build.Package, *go_ast.Package, error) {
	var (
		build_pkg *go_build.Package
		ast_pkg   *go_ast.Package
	)

	{ // lookup in memory
		pkg, err := ctx.go_ctx.Import(dir, srcdir, go_build.FindOnly)

		if _, ok := err.(*go_build.NoGoError); ok {
			err = nil
		}

		if err != nil {
			return nil, nil, err
		}

		if pkg, p := ctx.go_build_packages[pkg.ImportPath]; p {
			build_pkg = pkg
		}
	}

	if build_pkg == nil { // lookup in fs
		pkg, err := ctx.go_ctx.Import(dir, srcdir, 0)

		if _, ok := err.(*go_build.NoGoError); ok {
			err = nil
		}

		if err != nil {
			return nil, nil, err
		}

		ctx.go_build_packages[pkg.ImportPath] = pkg
		build_pkg = pkg
	}

	{ // lookup in mem
		if pkg, p := ctx.go_ast_packages[build_pkg.ImportPath]; p {
			ast_pkg = pkg
		}
	}

	if ast_pkg == nil { // lookup in fs
		ast_pkgs, err := go_parser.ParseDir(
			ctx.go_fset,
			build_pkg.Dir,
			go_file_filter(build_pkg),
			go_parser.SpuriousErrors,
		)
		if err != nil {
			return nil, nil, err
		}

		// collect files
		files := map[string]*go_ast.File{}
		for _, p := range ast_pkgs {
			for n, f := range p.Files {
				files[n] = f
			}
		}

		// make new package
		importer := func(imports map[string]*go_ast.Object, dir string) (*go_ast.Object, error) {
			return ctx.go_package_importer(imports, dir, build_pkg.ImportPath)
		}
		pkg, err := go_ast.NewPackage(ctx.go_fset, files, importer, Universe)
		if err != nil {
			return nil, nil, err
		}

		ctx.go_ast_packages[build_pkg.ImportPath] = pkg
		ast_pkg = pkg
	}

	return build_pkg, ast_pkg, nil
}

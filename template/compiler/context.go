package compiler

import (
	"fmt"
	w_ast "github.com/fd/w/template/ast"
	go_ast "go/ast"
	go_build "go/build"
	go_parser "go/parser"
	go_token "go/token"
	"os"
	"strings"
)

type Context struct {
	WROOT string

	DataViews   map[string]*DataView
	RenderFuncs map[string]*RenderFunc
	Helpers     map[string]*go_ast.FuncDecl

	go_ctx            *go_build.Context
	go_fset           *go_token.FileSet
	go_ast_packages   map[string]*go_ast.Package
	go_build_packages map[string]*go_build.Package
}

type RenderFunc struct {
	Name       string
	ImportPath string
	Template   *w_ast.Template
	Export     bool
	Golang     string
}

func (f *RenderFunc) FunctionName() string {
	if f.Export == false {
		return f.Name
	}

	parts := strings.Split(f.Name, "_")
	for i, part := range parts {
		parts[i] = strings.Title(part)
	}
	return strings.Join(parts, "")
}

type DataView struct {
	Name       string
	ImportPath string
	Expression w_ast.Expression
}

type Include struct {
	w_ast.Info
	View *DataView
}

func (n *Include) String() string {
	return "{{include}}"
}

type Errors []error

func (errs Errors) Error() string {
	s := "Errors:\n"
	for _, err := range errs {
		s += " - " + err.Error() + "\n"
	}
	return s
}

func (n *Include) Visit(b w_ast.Visitor) {
	// do nothing
}

func (ctx *Context) Analyze(dir string) error {
	if ctx.go_ctx == nil {
		// copy the context
		c := go_build.Default
		ctx.go_ctx = &c
		ctx.go_ctx.CgoEnabled = false
	}

	if ctx.go_fset == nil {
		ctx.go_fset = go_token.NewFileSet()
	}

	if ctx.go_ast_packages == nil {
		ctx.go_ast_packages = make(map[string]*go_ast.Package)
	}

	if ctx.go_build_packages == nil {
		ctx.go_build_packages = make(map[string]*go_build.Package)
	}

	if ctx.Helpers == nil {
		ctx.Helpers = make(map[string]*go_ast.FuncDecl)
	}

	if ctx.RenderFuncs == nil {
		ctx.RenderFuncs = make(map[string]*RenderFunc)
	}

	if ctx.DataViews == nil {
		ctx.DataViews = make(map[string]*DataView)
	}

	build_pkg, ast_pkg, err := ctx.ParsePackage(dir)
	if err != nil {
		return err
	}

	if build_pkg != nil && ast_pkg != nil {
		ctx.go_ast_packages[build_pkg.ImportPath] = ast_pkg
	}

	if build_pkg != nil {
		ctx.go_build_packages[build_pkg.ImportPath] = build_pkg
	}

	return nil
}

func (ctx *Context) ParsePackage(dir string) (*go_build.Package, *go_ast.Package, error) {
	pkg, err := ctx.go_ctx.Import(dir, ctx.WROOT, 0)
	if err != nil {
		if strings.HasPrefix(err.Error(), "no Go source files in") {
			err = nil
			return pkg, nil, nil
		}
		return nil, nil, err
	}

	pkgs, err := go_parser.ParseDir(ctx.go_fset, pkg.Dir, go_file_filter(pkg), go_parser.SpuriousErrors)
	if err != nil {
		return nil, nil, err
	}

	for n, p := range pkgs {
		if n == pkg.Name {
			return pkg, p, nil
		}
	}

	return nil, nil, fmt.Errorf("package not found: %s", dir)
}

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

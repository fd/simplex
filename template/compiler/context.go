package compiler

import (
	"fmt"
	w_ast "github.com/fd/w/template/ast"
	w_parser "github.com/fd/w/template/parser"
	go_ast "go/ast"
	go_build "go/build"
	go_parser "go/parser"
	go_token "go/token"
	"os"
	"path"
	"strings"
)

type Context struct {
	WROOT string

	DataViews   map[string]*DataView
	RenderFuncs map[string]*RenderFunc
	Helpers     map[string]*go_ast.FuncDecl

	go_ctx *go_build.Context
	fset   *go_token.FileSet
}

type RenderFunc struct {
	Name       string
	ImportPath string
	Template   *w_ast.Template
	Export     bool
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

func (ctx *Context) Analyze(dir string) error {
	if ctx.go_ctx == nil {
		// copy the context
		c := go_build.Default
		ctx.go_ctx = &c
		ctx.go_ctx.CgoEnabled = false
	}

	if ctx.fset == nil {
		ctx.fset = go_token.NewFileSet()
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

	err = ctx.FindFunctions(build_pkg.ImportPath, ast_pkg)
	if err != nil {
		return err
	}

	err = ctx.ParseTemplates(build_pkg.ImportPath, build_pkg.Dir)
	if err != nil {
		return err
	}

	return nil
}

func (ctx *Context) ParsePackage(dir string) (*go_build.Package, *go_ast.Package, error) {
	pkg, err := ctx.go_ctx.Import(dir, ctx.WROOT, 0)
	if err != nil {
		return nil, nil, err
	}

	pkgs, err := go_parser.ParseDir(ctx.fset, pkg.Dir, go_file_filter(pkg), go_parser.SpuriousErrors)
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

func (ctx *Context) FindFunctions(path string, pkg *go_ast.Package) error {
	finder := &FuncFinder{
		Helpers:      ctx.Helpers,
		package_name: path,
	}
	go_ast.Walk(finder, pkg)
	return nil
}

type FuncFinder struct {
	Helpers map[string]*go_ast.FuncDecl

	package_name string
}

func (v *FuncFinder) Visit(n go_ast.Node) go_ast.Visitor {
	if f, ok := n.(*go_ast.FuncDecl); ok {
		v.AnalyzeFunc(f)
		return nil
	}
	return v
}

func (v *FuncFinder) AnalyzeFunc(f *go_ast.FuncDecl) {

	// helpers must have no receiver
	if f.Recv != nil {
		return
	}

	// helpers must be exported
	if !f.Name.IsExported() {
		return
	}

	fullname := fmt.Sprintf("\"%s\".%s", v.package_name, f.Name.String())
	v.Helpers[fullname] = f
}

func (ctx *Context) ParseTemplates(import_path, dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}

	entries, err := d.Readdir(0)
	if err != nil {
		return err
	}

	for _, fi := range entries {
		base := fi.Name()

		if !strings.HasSuffix(base, ".go.html") {
			continue
		}

		tmpl, err := w_parser.ParseFile(path.Join(dir, base))
		if err != nil {
			return err
		}

		base = base[:len(base)-8]

		name := fmt.Sprintf("\"%s\".%s", import_path, base)
		ctx.RenderFuncs[name] = &RenderFunc{
			Name:       base,
			ImportPath: import_path,
			Template:   tmpl,
			Export:     true,
		}
	}

	return nil
}

package compiler

import (
	"fmt"
	w_ast "github.com/fd/w/template/ast"
	go_ast "go/ast"
	go_build "go/build"
	go_token "go/token"
	"strings"
)

type Context struct {
	WROOT string

	DataViews   map[string]*DataView
	RenderFuncs map[string]*RenderFunc
	Helpers     map[string]*go_ast.FuncDecl
	Imports     map[string]*Imports

	go_ctx            *go_build.Context
	go_fset           *go_token.FileSet
	go_universe       *go_ast.Scope
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

func (n *Include) Visit(b w_ast.Visitor) {
	// do nothing
}

type Errors []error

func (errs Errors) Error() string {
	s := "Errors:"
	for _, err := range errs {
		s += "\n - " + err.Error()
	}
	return s
}

func (errs Errors) Any() error {
	if len(errs) == 0 {
		return nil
	}
	return errs
}

type Imports struct {
	next_id int
	imports map[string]string
}

func (i *Imports) Register(import_path string) string {
	if i.imports == nil {
		i.imports = map[string]string{}
	}

	name, p := i.imports[import_path]
	if !p {
		i.next_id += 1
		name = fmt.Sprintf("import_%d", i.next_id)
		i.imports[import_path] = name
	}
	return name
}

func NewContext(wroot string) *Context {
	ctx := &Context{WROOT: wroot}

	if ctx.go_ctx == nil {
		// copy the context
		c := go_build.Default
		ctx.go_ctx = &c
		ctx.go_ctx.CgoEnabled = false
	}

	ctx.go_fset = go_token.NewFileSet()
	ctx.go_universe = go_ast.NewScope(nil)
	ctx.go_ast_packages = make(map[string]*go_ast.Package)
	ctx.go_build_packages = make(map[string]*go_build.Package)
	ctx.Helpers = make(map[string]*go_ast.FuncDecl)
	ctx.RenderFuncs = make(map[string]*RenderFunc)
	ctx.DataViews = make(map[string]*DataView)
	ctx.Imports = make(map[string]*Imports)

	return ctx
}

func (ctx *Context) Compile() error {
	var err error

	err = ctx.ImportPackages()
	if err != nil {
		return err
	}

	ctx.GolangFindFunctions()

	err = ctx.ParseTemplates()
	if err != nil {
		return err
	}

	ctx.LookupFunctionCalls()

	ctx.NormalizeGetExpresions()

	ctx.UnfoldRenderFunctions()

	ctx.CleanTemplates()

	ctx.PrintRenderFunctions()

	return nil
}

func (ctx *Context) ImportsFor(pkg string) *Imports {
	i, p := ctx.Imports[pkg]
	if !p {
		i = &Imports{}
		ctx.Imports[pkg] = i
	}
	return i
}

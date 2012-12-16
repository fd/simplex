package compiler

import (
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

	return ctx
}

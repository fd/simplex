package compiler

import (
	sx_ast "github.com/fd/w/simplex/ast"
	go_ast "go/ast"
	"go/build"
	"go/token"
	"time"
)

var Context build.Context
var Packages = map[string]*Package{}

func init() {
	Context = build.Default
	Context.CgoEnabled = false
}

type Package struct {
	TargetPath         string
	Imports            map[string]*Package
	SmplxTemplateFiles []string
	BuildPackage       *build.Package
	FileSet            *token.FileSet
	Files              map[string]*go_ast.File
	SmplxFiles         map[string]*sx_ast.File
	ModTimes           map[string]time.Time
	AstPackage         *go_ast.Package
	Views              map[string]*ViewDecl

	GeneratedFile *go_ast.File
}

type ViewDecl struct {
	MemberType *go_ast.Ident
	ViewType   *go_ast.Ident

	Source        *go_ast.FuncDecl
	Select        *go_ast.FuncDecl
	Reject        *go_ast.FuncDecl
	Sort          *go_ast.FuncDecl
	Group         *go_ast.FuncDecl
	CollectedFrom *go_ast.FuncDecl
}

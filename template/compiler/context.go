package compiler

import (
	"fmt"
	go_ast "go/ast"
	go_parser "go/parser"
	go_token "go/token"
	"os"
	"path"
)

func Context(dir string) error {
	fset := go_token.NewFileSet()

	filter := func(f os.FileInfo) bool {
		if path.Ext(f.Name()) != ".go" {
			return false
		}

		if f.IsDir() {
			return false
		}

		return true
	}

	pkgs, err := go_parser.ParseDir(fset, dir, filter, go_parser.SpuriousErrors)
	if err != nil {
		return err
	}

	for pkg_name, pkg := range pkgs {
		finder := &FuncFinder{
			Helpers: make(map[string]*go_ast.FuncDecl),
		}

		go_ast.Walk(finder, pkg)

		for n, f := range finder.Helpers {
			fmt.Printf("======> %s.%s\n", pkg_name, n)
			go_ast.Print(fset, f)
		}
	}

	return nil
}

type FuncFinder struct {
	Helpers map[string]*go_ast.FuncDecl
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

	v.Helpers[f.Name.String()] = f

}

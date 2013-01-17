package compiler

import (
	"bytes"
	"fmt"
	"github.com/fd/w/simplex/ast"
	"github.com/fd/w/simplex/token"
)

type Context struct {
	OutputFile string
	GoFiles    []string
	SxFiles    []string

	AstFiles   map[string]*ast.File
	AstPackage *ast.Package
	FileSet    *token.FileSet
}

func (c *Context) Compile() error {
	var err error

	err = c.parse_files()
	if err != nil {
		return err
	}

	err = c.check_types()
	if err != nil {
		return err
	}

	return nil
}

func (c *Context) generate_go() error {
	var b bytes.Buffer

	fmt.Fprintf(&b, intro, c.AstPackage.Name)

	// print table interfaces
	// print table structs
	// print keyed view interfaces
	// print keyed view structs
	// print indexed view interfaces
	// print indexed view structs
	ast.Walk(&type_printer{&b}, c.AstPackage)

	// parse generated go file
	// replace type expr in .sx files
	// replace view methods in .sx files
	// merge .sx files into generated go file
	// print generated go file

	return nil
}

const intro = `
package %s

import (
	sx_runtime "github.com/fd/w/simplex/runtime"
)
`

type type_printer struct {
	b *bytes.Buffer
}

func (v *type_printer) Visit(n ast.Node) ast.Visitor {
	return nil
}

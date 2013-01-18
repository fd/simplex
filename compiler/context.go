package compiler

import (
	"bytes"
	"fmt"
	"github.com/fd/simplex/ast"
	"github.com/fd/simplex/token"
	"github.com/fd/simplex/types"
)

type Context struct {
	OutputFile string
	GoFiles    []string
	SxFiles    []string

	AstFiles     map[string]*ast.File
	TypesPackage *types.Package
	ViewTypes    map[string]*types.View
	TableTypes   map[string]*types.Table
	FileSet      *token.FileSet
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

	err = c.generate_go()
	if err != nil {
		return err
	}

	return nil
}

func (c *Context) generate_go() error {
	var b bytes.Buffer

	fmt.Fprintf(&b, intro, c.TypesPackage.Name)

	for view := range c.ViewTypes {
		fmt.Printf("- view: %+v\n", view)
	}
	for table := range c.TableTypes {
		fmt.Printf("- table: %+v\n", table)
	}

	// print table interfaces
	// print table structs
	// print keyed view interfaces
	// print keyed view structs
	// print indexed view interfaces
	// print indexed view structs
	//ast.Walk(&type_printer{&b}, c.AstPackage)

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
	sx_runtime "github.com/fd/simplex/runtime"
)
`

type type_printer struct {
	b *bytes.Buffer
}

func (v *type_printer) Visit(n ast.Node) ast.Visitor {
	return nil
}

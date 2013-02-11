package compiler

import (
	"github.com/fd/simplex/lang/ast"
	"github.com/fd/simplex/lang/token"
	"github.com/fd/simplex/lang/types"
)

type Context struct {
	OutputDir  string
	ImportPath string
	GoFiles    []string
	SxFiles    []string

	AstFiles     map[string]*ast.File
	TypesPackage *types.Package
	ViewTypes    map[string]*types.View
	TableTypes   map[string]*types.Table
	FileSet      *token.FileSet

	NodeTypes map[ast.Node]types.Type
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

	err = c.convert_sx_to_go()
	if err != nil {
		return err
	}

	err = c.print_go()
	if err != nil {
		return err
	}

	return nil
}

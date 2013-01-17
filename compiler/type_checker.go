package compiler

import (
	"github.com/fd/simplex/ast"
	"github.com/fd/simplex/types"
)

func (c *Context) check_types() error {
	files := make([]*ast.File, 0, len(c.AstFiles))

	for _, name := range c.GoFiles {
		files = append(files, c.AstFiles[name])
	}

	for _, name := range c.SxFiles {
		files = append(files, c.AstFiles[name])
	}

	pkg, err := types.Check(c.FileSet, files)
	if err != nil {
		return err
	}

	c.TypesPackage = pkg
	return nil
}

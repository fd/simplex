package compiler

import (
	"github.com/fd/simplex/types"
)

func (c *Context) check_types() error {
	pkg, err := types.Check(c.FileSet, c.AstFiles)
	if err != nil {
		return err
	}

	c.AstPackage = pkg
	return nil
}

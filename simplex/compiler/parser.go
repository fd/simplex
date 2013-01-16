package compiler

import (
	"github.com/fd/w/simplex/ast"
	"github.com/fd/w/simplex/parser"
	"github.com/fd/w/simplex/token"
)

func (c *Context) parse_files() error {
	c.FileSet = token.NewFileSet()
	c.AstFiles = make(map[string]*ast.File, len(c.GoFiles)+len(c.SxFiles))

	for _, name := range c.GoFiles {
		file, err := parser.ParseFile(c.FileSet, name, nil, parser.ParseComments|parser.SpuriousErrors)
		if err != nil {
			return err
		}
		c.AstFiles[name] = file
	}

	for _, name := range c.SxFiles {
		file, err := parser.ParseFile(c.FileSet, name, nil, parser.ParseComments|parser.SpuriousErrors|parser.SimplexExtentions)
		if err != nil {
			return err
		}
		c.AstFiles[name] = file
	}

	return nil
}

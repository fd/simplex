package compiler

import (
	"fmt"
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

	views := map[*types.View]bool{}
	tables := map[*types.Table]bool{}
	cache := map[types.Type]bool{}

	ctx := types.Default
	ctx.Expr = func(x ast.Expr, typ types.Type, val interface{}) {
		collect_types(typ, views, tables, cache)
	}

	pkg, err := ctx.Check(c.FileSet, files)
	if err != nil {
		return err
	}

	c.TypesPackage = pkg

	for view, _ := range views {
		fmt.Printf("- view: %+v\n", view)
	}
	for table, _ := range tables {
		fmt.Printf("- table: %+v\n", table)
	}

	return nil
}

func collect_types(typ types.Type, views map[*types.View]bool, tables map[*types.Table]bool, cache map[types.Type]bool) {
	if _, p := cache[typ]; p {
		return
	}
	cache[typ] = true

	switch t := typ.(type) {
	case *types.View:
		if _, p := views[t]; p {
			return
		}

		views[t] = true
		collect_types(t.Elt, views, tables, cache)
	case *types.Table:
		if _, p := tables[t]; p {
			return
		}

		tables[t] = true
		collect_types(t.Elt, views, tables, cache)

	case *types.Array:
		collect_types(t.Elt, views, tables, cache)

	case *types.Slice:
		collect_types(t.Elt, views, tables, cache)

	case *types.Struct:
		for _, f := range t.Fields {
			collect_types(f.Type, views, tables, cache)
		}

	case *types.Pointer:
		collect_types(t.Base, views, tables, cache)

	case *types.Result:
		for _, p := range t.Values {
			collect_types(p.Type, views, tables, cache)
		}

	case *types.Signature:
		if t.Recv != nil {
			collect_types(t.Recv.Type, views, tables, cache)
		}
		for _, p := range t.Params {
			collect_types(p.Type, views, tables, cache)
		}
		for _, p := range t.Results {
			collect_types(p.Type, views, tables, cache)
		}

	case *types.Interface:
		for _, m := range t.Methods {
			collect_types(m.Type, views, tables, cache)
		}

	case *types.Map:
		collect_types(t.Key, views, tables, cache)
		collect_types(t.Elt, views, tables, cache)

	case *types.Chan:
		collect_types(t.Elt, views, tables, cache)

	case *types.NamedType:
		collect_types(t.Underlying, views, tables, cache)
		for _, m := range t.Methods {
			collect_types(m.Type, views, tables, cache)
		}

	}
}

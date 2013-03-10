package compiler

import (
	"simplex.sh/lang/ast"
	"simplex.sh/lang/types"
)

func (c *Context) check_types() error {
	files := make([]*ast.File, 0, len(c.AstFiles))

	for _, name := range c.GoFiles {
		files = append(files, c.AstFiles[name])
	}

	for _, name := range c.SxFiles {
		files = append(files, c.AstFiles[name])
	}

	views := map[string]*types.View{}
	tables := map[string]*types.Table{}
	cache := map[types.Type]bool{}
	mapping := map[ast.Node]types.Type{}

	ctx := types.Context{}
	ctx.Expr = func(x ast.Expr, typ types.Type, val interface{}) {
		mapping[x] = typ
		collect_types(typ, views, tables, cache)
	}

	pkg, err := ctx.Check(c.FileSet, files)
	if err != nil {
		return err
	}

	c.TypesPackage = pkg
	c.ViewTypes = views
	c.TableTypes = tables
	c.NodeTypes = mapping

	return nil
}

func collect_types(typ types.Type, views map[string]*types.View, tables map[string]*types.Table, cache map[types.Type]bool) {
	if _, p := cache[typ]; p {
		return
	}
	cache[typ] = true

	switch t := typ.(type) {
	case *types.View:
		name := view_type_name(t)
		views[name] = t
		collect_types(t.Elt, views, tables, cache)
	case *types.Table:
		name := view_type_name(t)
		tables[name] = t
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

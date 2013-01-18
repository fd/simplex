package compiler

import (
	"bytes"
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

	views := map[string]*types.View{}
	tables := map[string]*types.Table{}
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
	c.ViewTypes = views
	c.TableTypes = tables

	return nil
}

func collect_types(typ types.Type, views map[string]*types.View, tables map[string]*types.Table, cache map[types.Type]bool) {
	if _, p := cache[typ]; p {
		return
	}
	cache[typ] = true

	switch t := typ.(type) {
	case *types.View:
		name := type_name(t)
		views[name] = t
		collect_types(t.Elt, views, tables, cache)
	case *types.Table:
		name := type_name(t)
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

func type_name(typ types.Type) string {
	var b bytes.Buffer
	write_type_name(&b, typ)
	return b.String()
}

func write_type_name(b *bytes.Buffer, typ types.Type) {
	switch t := typ.(type) {
	case *types.View:
		b.WriteString("SxVi")
		if t.Key != nil {
			b.WriteString("_k")
			write_type_name(b, t.Key)
		}
		b.WriteString("_v")
		write_type_name(b, t.Elt)

	case *types.Table:
		b.WriteString("SxTa")
		b.WriteString("_k")
		write_type_name(b, t.Key)
		b.WriteString("_v")
		write_type_name(b, t.Elt)

	case *types.Basic:
		switch t.Kind {
		case types.Bool:
			b.WriteString("bool")
		case types.Int:
			b.WriteString("int")
		case types.Int8:
			b.WriteString("int8")
		case types.Int16:
			b.WriteString("int16")
		case types.Int32:
			b.WriteString("int32")
		case types.Int64:
			b.WriteString("int64")
		case types.Uint:
			b.WriteString("uint")
		case types.Uint8:
			b.WriteString("uint8")
		case types.Uint16:
			b.WriteString("uint16")
		case types.Uint32:
			b.WriteString("uint32")
		case types.Uint64:
			b.WriteString("uint64")
		case types.Uintptr:
			b.WriteString("uintptr")
		case types.Float32:
			b.WriteString("float32")
		case types.Float64:
			b.WriteString("float64")
		case types.Complex64:
			b.WriteString("complex64")
		case types.Complex128:
			b.WriteString("complex128")
		case types.String:
			b.WriteString("string")
		case types.UnsafePointer:
			b.WriteString("UnsafePointer")
		}

	case *types.Array:
		fmt.Fprintf(b, "Ar%d_", t.Len)
		write_type_name(b, t.Elt)

	case *types.Slice:
		b.WriteString("Sl_")
		write_type_name(b, t.Elt)

	case *types.Struct:
		b.WriteString("St")
		for _, f := range t.Fields {
			fmt.Fprintf(b, "_n%s_", f.Name)
			write_type_name(b, f.Type)
		}

	case *types.Pointer:
		b.WriteString("Pt_")
		write_type_name(b, t.Base)

	case *types.Result:
		panic("not identifiable")

	case *types.Signature:
		b.WriteString("Si")
		if t.Recv != nil {
			b.WriteString("_r")
			write_type_name(b, t.Recv.Type)
		}
		for _, p := range t.Params {
			b.WriteString("_i")
			write_type_name(b, p.Type)
		}
		for _, p := range t.Results {
			b.WriteString("_o")
			write_type_name(b, p.Type)
		}

	case *types.Interface:
		b.WriteString("In")
		for _, f := range t.Methods {
			fmt.Fprintf(b, "_n%s_", f.Name)
			write_type_name(b, f.Type)
		}

	case *types.Map:
		b.WriteString("Ma")
		b.WriteString("_k")
		write_type_name(b, t.Key)
		b.WriteString("_v")
		write_type_name(b, t.Elt)

	case *types.Chan:
		b.WriteString("Ch_")
		write_type_name(b, t.Elt)

	case *types.NamedType:
		b.WriteString(t.Obj.Name)

	default:
		panic("unhandle type name")
	}
}

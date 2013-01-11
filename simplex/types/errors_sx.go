package types

import (
	"bytes"
	"fmt"
	"github.com/fd/w/simplex/ast"
)

// TODO(gri) Need to merge with typeString since some expressions are types (try: ([]int)(a))
//
// see errors.go:115
func writeExpr(buf *bytes.Buffer, expr ast.Expr) {
	switch x := expr.(type) {
	case *ast.Ident:
		buf.WriteString(x.Name)

	case *ast.BasicLit:
		buf.WriteString(x.Value)

	case *ast.FuncLit:
		buf.WriteString("(func literal)")

	case *ast.CompositeLit:
		buf.WriteString("(composite literal)")

	case *ast.ParenExpr:
		buf.WriteByte('(')
		writeExpr(buf, x.X)
		buf.WriteByte(')')

	case *ast.SelectorExpr:
		writeExpr(buf, x.X)
		buf.WriteByte('.')
		buf.WriteString(x.Sel.Name)

	case *ast.IndexExpr:
		writeExpr(buf, x.X)
		buf.WriteByte('[')
		writeExpr(buf, x.Index)
		buf.WriteByte(']')

	case *ast.SliceExpr:
		writeExpr(buf, x.X)
		buf.WriteByte('[')
		if x.Low != nil {
			writeExpr(buf, x.Low)
		}
		buf.WriteByte(':')
		if x.High != nil {
			writeExpr(buf, x.High)
		}
		buf.WriteByte(']')

	case *ast.TypeAssertExpr:
		writeExpr(buf, x.X)
		buf.WriteString(".(...)")

	case *ast.CallExpr:
		writeExpr(buf, x.Fun)
		buf.WriteByte('(')
		for i, arg := range x.Args {
			if i > 0 {
				buf.WriteString(", ")
			}
			writeExpr(buf, arg)
		}
		buf.WriteByte(')')

	case *ast.StarExpr:
		buf.WriteByte('*')
		writeExpr(buf, x.X)

	case *ast.UnaryExpr:
		buf.WriteString(x.Op.String())
		writeExpr(buf, x.X)

	case *ast.BinaryExpr:
		// The AST preserves source-level parentheses so there is
		// no need to introduce parentheses here for correctness.
		writeExpr(buf, x.X)
		buf.WriteByte(' ')
		buf.WriteString(x.Op.String())
		buf.WriteByte(' ')
		writeExpr(buf, x.Y)

	//=== start custom
	case *ast.StepExpr:
		writeExpr(buf, x.X)
		buf.WriteByte('.')
		buf.WriteString(x.StepType.String())
		buf.WriteByte('(')
		writeExpr(buf, x.F)
		buf.WriteByte(')')
	//=== end custom

	default:
		fmt.Fprintf(buf, "<expr %T>", x)
	}
}

// see errors.go:247
func writeType(buf *bytes.Buffer, typ Type) {
	switch t := typ.(type) {
	case nil:
		buf.WriteString("<nil>")

	case *Basic:
		buf.WriteString(t.Name)

	case *Array:
		fmt.Fprintf(buf, "[%d]", t.Len)
		writeType(buf, t.Elt)

	case *Slice:
		buf.WriteString("[]")
		writeType(buf, t.Elt)

	case *Struct:
		buf.WriteString("struct{")
		for i, f := range t.Fields {
			if i > 0 {
				buf.WriteString("; ")
			}
			if !f.IsAnonymous {
				buf.WriteString(f.Name)
				buf.WriteByte(' ')
			}
			writeType(buf, f.Type)
			if f.Tag != "" {
				fmt.Fprintf(buf, " %q", f.Tag)
			}
		}
		buf.WriteByte('}')

	case *Pointer:
		buf.WriteByte('*')
		writeType(buf, t.Base)

	case *Result:
		writeParams(buf, t.Values, false)

	case *Signature:
		buf.WriteString("func")
		writeSignature(buf, t)

	case *builtin:
		fmt.Fprintf(buf, "<type of %s>", t.name)

	case *Interface:
		buf.WriteString("interface{")
		for i, m := range t.Methods {
			if i > 0 {
				buf.WriteString("; ")
			}
			buf.WriteString(m.Name)
			writeSignature(buf, m.Type)
		}
		buf.WriteByte('}')

	case *Map:
		buf.WriteString("map[")
		writeType(buf, t.Key)
		buf.WriteByte(']')
		writeType(buf, t.Elt)

	case *Chan:
		var s string
		switch t.Dir {
		case ast.SEND:
			s = "chan<- "
		case ast.RECV:
			s = "<-chan "
		default:
			s = "chan "
		}
		buf.WriteString(s)
		writeType(buf, t.Elt)

	case *NamedType:
		buf.WriteString(t.Obj.Name)

	//=== start custom
	case *View:
		buf.WriteString("view[")
		if t.Key != nil {
			writeType(buf, t.Key)
		}
		buf.WriteByte(']')
		writeType(buf, t.Elt)

	case *Table:
		buf.WriteString("table[")
		writeType(buf, t.Key)
		buf.WriteByte(']')
		writeType(buf, t.Elt)
	//=== end custom

	default:
		fmt.Fprintf(buf, "<type %T>", t)
	}
}

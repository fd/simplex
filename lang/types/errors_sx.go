package types

import (
	"bytes"
	"fmt"
	"simplex.sh/lang/ast"
)

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
		s := "<NamedType w/o object>"
		if t.Obj != nil {
			s = t.Obj.GetName()
		}
		buf.WriteString(s)

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

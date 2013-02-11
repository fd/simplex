package compiler

import (
	"bytes"
	"fmt"
	"github.com/fd/simplex/lang/ast"
	"github.com/fd/simplex/lang/types"
)

func sx_type_string(typ types.Type) string {
	var b bytes.Buffer
	write_sx_type_string(&b, typ)
	return b.String()
}

func write_sx_type_string(b *bytes.Buffer, typ types.Type) {
	switch t := typ.(type) {
	case *types.View:
		b.WriteString("view[")
		if t.Key != nil {
			write_sx_type_string(b, t.Key)
		}
		b.WriteByte(']')
		write_sx_type_string(b, t.Elt)

	case *types.Table:
		b.WriteString("table[")
		write_sx_type_string(b, t.Key)
		b.WriteByte(']')
		write_sx_type_string(b, t.Elt)

	case *types.Basic:
		b.WriteString(t.Name)

	case *types.Array:
		fmt.Fprintf(b, "[%d]", t.Len)
		write_sx_type_string(b, t.Elt)

	case *types.Slice:
		b.WriteByte('[')
		b.WriteByte(']')
		write_sx_type_string(b, t.Elt)

	case *types.Struct:
		b.WriteString("struct{")
		for i, f := range t.Fields {
			if i > 0 {
				b.WriteString("; ")
			}
			if !f.IsAnonymous {
				b.WriteString(f.Name)
				b.WriteByte(' ')
			}
			write_sx_type_string(b, f.Type)
			if f.Tag != "" {
				fmt.Fprintf(b, " %q", f.Tag)
			}
		}
		b.WriteByte('}')

	case *types.Pointer:
		b.WriteByte('*')
		write_sx_type_string(b, t.Base)

	case *types.Result:
		write_sx_type_params_string(b, t.Values, false)

	case *types.Signature:
		b.WriteString("func")
		write_sx_type_signature_string(b, t)

	case *types.Interface:
		b.WriteString("interface{")
		for i, f := range t.Methods {
			if i > 0 {
				b.WriteString("; ")
			}
			b.WriteString(f.Name)
			write_sx_type_signature_string(b, f.Type)
		}
		b.WriteByte('}')

	case *types.Map:
		b.WriteString("map[")
		write_sx_type_string(b, t.Key)
		b.WriteByte(']')
		write_sx_type_string(b, t.Elt)

	case *types.Chan:
		var s string
		switch t.Dir {
		case ast.SEND:
			s = "chan<- "
		case ast.RECV:
			s = "<-chan "
		default:
			s = "chan "
		}
		b.WriteString(s)
		write_sx_type_string(b, t.Elt)

	case *types.NamedType:
		s := "<NamedType w/o object>"
		if t.Obj != nil {
			s = t.Obj.GetName()
		}
		b.WriteString(s)

	default:
		panic(fmt.Sprintf("unhandle type %T", t))
	}
}

func write_sx_type_signature_string(buf *bytes.Buffer, sig *types.Signature) {
	write_sx_type_params_string(buf, sig.Params, sig.IsVariadic)
	if len(sig.Results) == 0 {
		// no result
		return
	}

	buf.WriteByte(' ')
	if len(sig.Results) == 1 && sig.Results[0].Name == "" {
		// single unnamed result
		write_sx_type_string(buf, sig.Results[0].Type.(types.Type))
		return
	}

	// multiple or named result(s)
	write_sx_type_params_string(buf, sig.Results, false)
}

func write_sx_type_params_string(buf *bytes.Buffer, params []*types.Var, isVariadic bool) {
	buf.WriteByte('(')
	for i, par := range params {
		if i > 0 {
			buf.WriteString(", ")
		}
		if par.Name != "" {
			buf.WriteString(par.Name)
			buf.WriteByte(' ')
		}
		if isVariadic && i == len(params)-1 {
			buf.WriteString("...")
		}
		write_sx_type_string(buf, par.Type)
	}
	buf.WriteByte(')')
}

func go_type_string(typ types.Type) string {
	var b bytes.Buffer
	write_go_type_string(&b, typ)
	return b.String()
}

func write_go_type_string(b *bytes.Buffer, typ types.Type) {
	switch t := typ.(type) {
	case *types.View:
		write_view_type_name(b, t)

	case *types.Table:
		write_view_type_name(b, t)

	case *types.Basic:
		b.WriteString(t.Name)

	case *types.Array:
		fmt.Fprintf(b, "[%d]", t.Len)
		write_go_type_string(b, t.Elt)

	case *types.Slice:
		b.WriteByte('[')
		b.WriteByte(']')
		write_go_type_string(b, t.Elt)

	case *types.Struct:
		b.WriteString("struct{")
		for i, f := range t.Fields {
			if i > 0 {
				b.WriteString("; ")
			}
			if !f.IsAnonymous {
				b.WriteString(f.Name)
				b.WriteByte(' ')
			}
			write_go_type_string(b, f.Type)
			if f.Tag != "" {
				fmt.Fprintf(b, " %q", f.Tag)
			}
		}
		b.WriteByte('}')

	case *types.Pointer:
		b.WriteByte('*')
		write_go_type_string(b, t.Base)

	case *types.Result:
		write_go_type_params_string(b, t.Values, false)

	case *types.Signature:
		b.WriteString("func")
		write_go_type_signature_string(b, t)

	case *types.Interface:
		b.WriteString("interface{")
		for i, f := range t.Methods {
			if i > 0 {
				b.WriteString("; ")
			}
			b.WriteString(f.Name)
			write_go_type_signature_string(b, f.Type)
		}
		b.WriteByte('}')

	case *types.Map:
		b.WriteString("map[")
		write_go_type_string(b, t.Key)
		b.WriteByte(']')
		write_go_type_string(b, t.Elt)

	case *types.Chan:
		var s string
		switch t.Dir {
		case ast.SEND:
			s = "chan<- "
		case ast.RECV:
			s = "<-chan "
		default:
			s = "chan "
		}
		b.WriteString(s)
		write_go_type_string(b, t.Elt)

	case *types.NamedType:
		s := "<NamedType w/o object>"
		if t.Obj != nil {
			s = t.Obj.GetName()
		}
		b.WriteString(s)

	default:
		panic(fmt.Sprintf("unhandle type %T", t))
	}
}

func write_go_type_signature_string(buf *bytes.Buffer, sig *types.Signature) {
	write_go_type_params_string(buf, sig.Params, sig.IsVariadic)
	if len(sig.Results) == 0 {
		// no result
		return
	}

	buf.WriteByte(' ')
	if len(sig.Results) == 1 && sig.Results[0].Name == "" {
		// single unnamed result
		write_go_type_string(buf, sig.Results[0].Type.(types.Type))
		return
	}

	// multiple or named result(s)
	write_go_type_params_string(buf, sig.Results, false)
}

func write_go_type_params_string(buf *bytes.Buffer, params []*types.Var, isVariadic bool) {
	buf.WriteByte('(')
	for i, par := range params {
		if i > 0 {
			buf.WriteString(", ")
		}
		if par.Name != "" {
			buf.WriteString(par.Name)
			buf.WriteByte(' ')
		}
		if isVariadic && i == len(params)-1 {
			buf.WriteString("...")
		}
		write_go_type_string(buf, par.Type)
	}
	buf.WriteByte(')')
}

/*
  print a zero type name for the views and tables
*/
func type_zero(typ types.Type) string {
	var b bytes.Buffer
	write_type_zero(&b, typ)
	return b.String()
}

func write_type_zero(b *bytes.Buffer, typ types.Type) {
	switch t := typ.(type) {
	case *types.View:
		b.WriteString("sx_")
		write_view_type_name(b, t)
		b.WriteString("{}")

	case *types.Table:
		b.WriteString("sx_")
		write_view_type_name(b, t)
		b.WriteString("{}")

	case *types.Basic:
		switch t.Kind {
		case types.Bool:
			b.WriteString("false")
		case types.Int, types.Int8, types.Int16, types.Int32, types.Int64:
			b.WriteString("0")
		case types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64:
			b.WriteString("0")
		case types.Uintptr:
			b.WriteString("0")
		case types.Float32, types.Float64:
			b.WriteString("0.0")
		case types.Complex64, types.Complex128:
			panic("no support for complex numbers")
		case types.String:
			b.WriteString("\"\"")
		case types.UnsafePointer:
			b.WriteString("0")
		}

	case *types.Array:
		b.WriteString("nil")

	case *types.Slice:
		b.WriteString("nil")

	case *types.Struct:
		write_go_type_string(b, t)
		b.WriteString("{}")

	case *types.Pointer:
		b.WriteString("nil")

	case *types.Result:
		panic("not zeroable")

	case *types.Signature:
		b.WriteString("nil")

	case *types.Interface:
		b.WriteString("nil")

	case *types.Map:
		b.WriteString("nil")

	case *types.Chan:
		b.WriteString("nil")

	case *types.NamedType:
		switch t.Underlying.(type) {

		case *types.Struct:
			b.WriteString(t.Obj.Name)
			b.WriteString("{}")

		case *types.Array, *types.Slice, *types.Pointer, *types.Signature, *types.Interface, *types.Map, *types.Chan:
			b.WriteString("nil")

		default:
			b.WriteString(t.Obj.Name)
			b.WriteByte('(')
			write_type_zero(b, t.Underlying)
			b.WriteByte(')')

		}

	default:
		panic("unhandle type name")
	}
}

/*
  print an identifiable type name for the views and tables
*/
func view_type_name(typ types.Type) string {
	var b bytes.Buffer
	write_view_type_name(&b, typ)
	return b.String()
}

func write_view_type_name(b *bytes.Buffer, typ types.Type) {
	switch t := typ.(type) {
	case *types.View:
		b.WriteString("SxVi")
		if t.Key != nil {
			b.WriteString("_k")
			write_view_type_name(b, t.Key)
		}
		b.WriteString("_v")
		write_view_type_name(b, t.Elt)

	case *types.Table:
		b.WriteString("SxTa")
		b.WriteString("_k")
		write_view_type_name(b, t.Key)
		b.WriteString("_v")
		write_view_type_name(b, t.Elt)

	case *types.Basic:
		b.WriteString(t.Name)

	case *types.Array:
		fmt.Fprintf(b, "Ar%d_", t.Len)
		write_view_type_name(b, t.Elt)

	case *types.Slice:
		b.WriteString("Sl_")
		write_view_type_name(b, t.Elt)

	case *types.Struct:
		b.WriteString("St")
		for _, f := range t.Fields {
			fmt.Fprintf(b, "_n%s_", f.Name)
			write_view_type_name(b, f.Type)
		}

	case *types.Pointer:
		b.WriteString("Pt_")
		write_view_type_name(b, t.Base)

	case *types.Result:
		panic("not identifiable")

	case *types.Signature:
		b.WriteString("Si")
		if t.Recv != nil {
			b.WriteString("_r")
			write_view_type_name(b, t.Recv.Type)
		}
		for _, p := range t.Params {
			b.WriteString("_i")
			write_view_type_name(b, p.Type)
		}
		for _, p := range t.Results {
			b.WriteString("_o")
			write_view_type_name(b, p.Type)
		}

	case *types.Interface:
		b.WriteString("In")
		for _, f := range t.Methods {
			fmt.Fprintf(b, "_n%s_", f.Name)
			write_view_type_name(b, f.Type)
		}

	case *types.Map:
		b.WriteString("Ma")
		b.WriteString("_k")
		write_view_type_name(b, t.Key)
		b.WriteString("_v")
		write_view_type_name(b, t.Elt)

	case *types.Chan:
		b.WriteString("Ch_")
		write_view_type_name(b, t.Elt)

	case *types.NamedType:
		b.WriteString(t.Obj.Name)

	default:
		panic("unhandle type name")
	}
}

package types

import (
	"github.com/fd/simplex/ast"
)

// TODO(gri) The functions operand.isAssignable, checker.convertUntyped,
//           checker.isRepresentable, and checker.assignOperand are
//           overlapping in functionality. Need to simplify and clean up.

// isAssignable reports whether x is assignable to a variable of type T.
//
// see operand.go:133
func (x *operand) isAssignable(T Type) bool {
	if x.mode == invalid || T == Typ[Invalid] {
		return true // avoid spurious errors
	}

	V := x.typ

	// x's type is identical to T
	if isIdentical(V, T) {
		return true
	}

	Vu := underlying(V)
	Tu := underlying(T)

	// x's type V and T have identical underlying types
	// and at least one of V or T is not a named type
	if isIdentical(Vu, Tu) {
		return !isNamed(V) || !isNamed(T)
	}

	// T is an interface type and x implements T
	if Ti, ok := Tu.(*Interface); ok {
		if m, _ := missingMethod(x.typ, Ti); m == nil {
			return true
		}
	}

	// x is a bidirectional channel value, T is a channel
	// type, x's type V and T have identical element types,
	// and at least one of V or T is not a named type
	if Vc, ok := Vu.(*Chan); ok && Vc.Dir == ast.SEND|ast.RECV {
		if Tc, ok := Tu.(*Chan); ok && isIdentical(Vc.Elt, Tc.Elt) {
			return !isNamed(V) || !isNamed(T)
		}
	}

	//=== start custom

	// x is a keyed view value, T is a view type,
	// x's type V and T have identical element types,
	// and at least one of V or T is not a named type
	if Vv, ok := Vu.(*View); ok {
		if Tv, ok := Tu.(*View); ok && Vv.Key != nil && isIdentical(Vv.Elt, Tv.Elt) {
			return !isNamed(V) || !isNamed(T)
		}
		if _, ok := Tu.(*Interface); ok && isNamed(T) {
			n := T.(*NamedType).Obj.GetName()
			if n == "GenericView" {
				return true
			}
		}
	}

	// x is a table value, T is a view keyed type,
	// x's type V and T have identical element types,
	// and at least one of V or T is not a named type
	if Vv, ok := Vu.(*Table); ok {
		if Tv, ok := Tu.(*View); ok && isIdentical(Vv.Elt, Tv.Elt) && isIdentical(Vv.Key, Tv.Key) {
			return !isNamed(V) || !isNamed(T)
		}
		if _, ok := Tu.(*Interface); ok && isNamed(T) {
			n := T.(*NamedType).Obj.GetName()
			if n == "GenericView" || n == "GenericTable" {
				return true
			}
		}
	}

	// x is a table value, T is a view type,
	// x's type V and T have identical element types,
	// and at least one of V or T is not a named type
	if Vv, ok := Vu.(*Table); ok {
		if Tv, ok := Tu.(*View); ok && isIdentical(Vv.Elt, Tv.Elt) && Tv.Key == nil {
			return !isNamed(V) || !isNamed(T)
		}
		if _, ok := Tu.(*Interface); ok && isNamed(T) {
			n := T.(*NamedType).Obj.GetName()
			if n == "GenericView" || n == "GenericTable" {
				return true
			}
		}
	}

	//=== end custom

	// x is the predeclared identifier nil and T is a pointer,
	// function, slice, map, channel, or interface type
	if x.isNil() {
		switch t := Tu.(type) {
		case *Basic:
			if t.Kind == UnsafePointer {
				return true
			}
		case *Pointer, *Signature, *Slice, *Map, *Chan, *Interface:
			return true

		//=== start custom
		case *View, *Table:
			return true
			//=== end custom
		}
		return false
	}

	// x is an untyped constant representable by a value of type T
	// TODO(gri) This is borrowing from checker.convertUntyped and
	//           checker.isRepresentable. Need to clean up.
	if isUntyped(Vu) {
		switch t := Tu.(type) {
		case *Basic:
			if x.mode == constant {
				return isRepresentableConst(x.val, t.Kind)
			}
			// The result of a comparison is an untyped boolean,
			// but may not be a constant.
			if Vb, _ := Vu.(*Basic); Vb != nil {
				return Vb.Kind == UntypedBool && isBoolean(Tu)
			}
		case *Interface:
			return x.isNil() || len(t.Methods) == 0
		case *Pointer, *Signature, *Slice, *Map, *Chan:
			return x.isNil()

		//=== start custom
		case *View, *Table:
			return x.isNil()
			//=== end custom

		}
	}

	return false
}

// see operand.go:354
func lookupField(typ Type, name QualifiedName) (operandMode, Type) {
	typ = deref(typ)

	if t, ok := typ.(*NamedType); ok {
		for _, m := range t.Methods {
			if name.IsSame(m.QualifiedName) {
				assert(m.Type != nil)
				return value, m.Type
			}
		}
		typ = t.Underlying
	}

	switch t := typ.(type) {
	case *Struct:
		var next []embeddedType
		for _, f := range t.Fields {
			if name.IsSame(f.QualifiedName) {
				return variable, f.Type
			}
			if f.IsAnonymous {
				// Possible optimization: If the embedded type
				// is a pointer to the current type we could
				// ignore it.
				next = append(next, embeddedType{typ: deref(f.Type).(*NamedType)})
			}
		}
		if len(next) > 0 {
			res := lookupFieldBreadthFirst(next, name)
			return res.mode, res.typ
		}

	case *Interface:
		for _, m := range t.Methods {
			if name.IsSame(m.QualifiedName) {
				return value, m.Type
			}
		}

	//=== start custom
	case *View, *Table:
		for _, n := range sx_step_names {
			if n == name.Name {
				return value, &Signature{Recv: &Var{Type: typ}}
			}
		}
		//=== end custom

	}

	// not found
	return invalid, nil
}

var sx_step_names = [...]string{
	"select",
	"reject",
	"detect",
	"collect",
	"inject",
	"group",
	"index",
	"sort",
}

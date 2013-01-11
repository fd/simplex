package types

// see predicates.go:88
func hasNil(typ Type) bool {
	switch underlying(typ).(type) {
	case *Slice, *Pointer, *Signature, *Interface, *Map, *Chan:
		return true

	//=== start custom
	case *View, *Table:
		return true
		//=== end custom

	}
	return false
}

// identical returns true if x and y are identical.
//
// see predicates.go:99
func isIdentical(x, y Type) bool {
	if x == y {
		return true
	}

	switch x := x.(type) {
	case *Basic:
		// Basic types are singletons except for the rune and byte
		// aliases, thus we cannot solely rely on the x == y check
		// above.
		if y, ok := y.(*Basic); ok {
			return x.Kind == y.Kind
		}

	case *Array:
		// Two array types are identical if they have identical element types
		// and the same array length.
		if y, ok := y.(*Array); ok {
			return x.Len == y.Len && isIdentical(x.Elt, y.Elt)
		}

	case *Slice:
		// Two slice types are identical if they have identical element types.
		if y, ok := y.(*Slice); ok {
			return isIdentical(x.Elt, y.Elt)
		}

	case *Struct:
		// Two struct types are identical if they have the same sequence of fields,
		// and if corresponding fields have the same names, and identical types,
		// and identical tags. Two anonymous fields are considered to have the same
		// name. Lower-case field names from different packages are always different.
		if y, ok := y.(*Struct); ok {
			// TODO(gri) handle structs from different packages
			if len(x.Fields) == len(y.Fields) {
				for i, f := range x.Fields {
					g := y.Fields[i]
					if f.Name != g.Name ||
						!isIdentical(f.Type, g.Type) ||
						f.Tag != g.Tag ||
						f.IsAnonymous != g.IsAnonymous {
						return false
					}
				}
				return true
			}
		}

	case *Pointer:
		// Two pointer types are identical if they have identical base types.
		if y, ok := y.(*Pointer); ok {
			return isIdentical(x.Base, y.Base)
		}

	case *Signature:
		// Two function types are identical if they have the same number of parameters
		// and result values, corresponding parameter and result types are identical,
		// and either both functions are variadic or neither is. Parameter and result
		// names are not required to match.
		if y, ok := y.(*Signature); ok {
			return identicalTypes(x.Params, y.Params) &&
				identicalTypes(x.Results, y.Results) &&
				x.IsVariadic == y.IsVariadic
		}

	case *Interface:
		// Two interface types are identical if they have the same set of methods with
		// the same names and identical function types. Lower-case method names from
		// different packages are always different. The order of the methods is irrelevant.
		if y, ok := y.(*Interface); ok {
			return identicalMethods(x.Methods, y.Methods) // methods are sorted
		}

	case *Map:
		// Two map types are identical if they have identical key and value types.
		if y, ok := y.(*Map); ok {
			return isIdentical(x.Key, y.Key) && isIdentical(x.Elt, y.Elt)
		}

	case *Chan:
		// Two channel types are identical if they have identical value types
		// and the same direction.
		if y, ok := y.(*Chan); ok {
			return x.Dir == y.Dir && isIdentical(x.Elt, y.Elt)
		}

	case *NamedType:
		// Two named types are identical if their type names originate
		// in the same type declaration.
		if y, ok := y.(*NamedType); ok {
			return x.Obj == y.Obj
		}

	//=== start custom
	case *View:
		// Two map types are identical if they have identical key and value types.
		if y, ok := y.(*View); ok {
			return isIdentical(x.Key, y.Key) && isIdentical(x.Elt, y.Elt)
		}

	case *Table:
		// Two map types are identical if they have identical key and value types.
		if y, ok := y.(*Table); ok {
			return isIdentical(x.Key, y.Key) && isIdentical(x.Elt, y.Elt)
		}
		//=== end custom

	}

	return false
}

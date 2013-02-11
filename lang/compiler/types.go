package compiler

import (
	"simplex.sh/lang/types"
)

func underlying_type(typ types.Type) types.Type {
	if named, ok := typ.(*types.NamedType); ok {
		return named.Underlying
	}
	return typ
}

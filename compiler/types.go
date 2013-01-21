package compiler

import (
	"github.com/fd/simplex/types"
)

func underlying_type(typ types.Type) types.Type {
	if named, ok := typ.(*types.NamedType); ok {
		return named.Underlying
	}
	return typ
}

// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file implements the universe and unsafe package scopes.

package types

import (
	"github.com/fd/w/simplex/ast"
	go_ast "go/ast"
	"strings"
)

var (
	aType            implementsType
	Universe, unsafe *go_ast.Scope
	Unsafe           *go_ast.Object // package unsafe
)

// Predeclared types, indexed by BasicKind.
var Typ = [...]*Basic{
	Invalid: {aType, Invalid, 0, 0, "invalid type"},

	Bool:          {aType, Bool, IsBoolean, 1, "bool"},
	Int:           {aType, Int, IsInteger, 0, "int"},
	Int8:          {aType, Int8, IsInteger, 1, "int8"},
	Int16:         {aType, Int16, IsInteger, 2, "int16"},
	Int32:         {aType, Int32, IsInteger, 4, "int32"},
	Int64:         {aType, Int64, IsInteger, 8, "int64"},
	Uint:          {aType, Uint, IsInteger | IsUnsigned, 0, "uint"},
	Uint8:         {aType, Uint8, IsInteger | IsUnsigned, 1, "uint8"},
	Uint16:        {aType, Uint16, IsInteger | IsUnsigned, 2, "uint16"},
	Uint32:        {aType, Uint32, IsInteger | IsUnsigned, 4, "uint32"},
	Uint64:        {aType, Uint64, IsInteger | IsUnsigned, 8, "uint64"},
	Uintptr:       {aType, Uintptr, IsInteger | IsUnsigned, 0, "uintptr"},
	Float32:       {aType, Float32, IsFloat, 4, "float32"},
	Float64:       {aType, Float64, IsFloat, 8, "float64"},
	Complex64:     {aType, Complex64, IsComplex, 8, "complex64"},
	Complex128:    {aType, Complex128, IsComplex, 16, "complex128"},
	String:        {aType, String, IsString, 0, "string"},
	UnsafePointer: {aType, UnsafePointer, 0, 0, "Pointer"},

	UntypedBool:    {aType, UntypedBool, IsBoolean | IsUntyped, 0, "untyped boolean"},
	UntypedInt:     {aType, UntypedInt, IsInteger | IsUntyped, 0, "untyped integer"},
	UntypedRune:    {aType, UntypedRune, IsInteger | IsUntyped, 0, "untyped rune"},
	UntypedFloat:   {aType, UntypedFloat, IsFloat | IsUntyped, 0, "untyped float"},
	UntypedComplex: {aType, UntypedComplex, IsComplex | IsUntyped, 0, "untyped complex"},
	UntypedString:  {aType, UntypedString, IsString | IsUntyped, 0, "untyped string"},
	UntypedNil:     {aType, UntypedNil, IsUntyped, 0, "untyped nil"},
}

var aliases = [...]*Basic{
	{aType, Byte, IsInteger | IsUnsigned, 1, "byte"},
	{aType, Rune, IsInteger, 4, "rune"},
}

var predeclaredConstants = [...]*struct {
	kind BasicKind
	name string
	val  interface{}
}{
	{UntypedBool, "true", true},
	{UntypedBool, "false", false},
	{UntypedInt, "iota", zeroConst},
	{UntypedNil, "nil", nilConst},
}

var predeclaredFunctions = [...]*builtin{
	{aType, _Append, "append", 1, true, false},
	{aType, _Cap, "cap", 1, false, false},
	{aType, _Close, "close", 1, false, true},
	{aType, _Complex, "complex", 2, false, false},
	{aType, _Copy, "copy", 2, false, true},
	{aType, _Delete, "delete", 2, false, true},
	{aType, _Imag, "imag", 1, false, false},
	{aType, _Len, "len", 1, false, false},
	{aType, _Make, "make", 1, true, false},
	{aType, _New, "new", 1, false, false},
	{aType, _Panic, "panic", 1, false, true},
	{aType, _Print, "print", 1, true, true},
	{aType, _Println, "println", 1, true, true},
	{aType, _Real, "real", 1, false, false},
	{aType, _Recover, "recover", 0, false, true},

	{aType, _Alignof, "Alignof", 1, false, false},
	{aType, _Offsetof, "Offsetof", 1, false, false},
	{aType, _Sizeof, "Sizeof", 1, false, false},
}

// commonly used types
var (
	emptyInterface = new(Interface)
)

// commonly used constants
var (
	universeIota *go_ast.Object
)

func init() {
	// Universe scope
	Universe = go_ast.NewScope(nil)

	// unsafe package and its scope
	unsafe = go_ast.NewScope(nil)
	Unsafe = go_ast.NewObj(go_ast.Pkg, "unsafe")
	Unsafe.Data = unsafe

	// predeclared types
	for _, t := range Typ {
		def(go_ast.Typ, t.Name).Type = t
	}
	for _, t := range aliases {
		def(go_ast.Typ, t.Name).Type = t
	}

	// error type
	{
		err := &Method{"Error", &Signature{Results: []*Var{{"", Typ[String]}}}}
		obj := def(go_ast.Typ, "error")
		obj.Type = &NamedType{Underlying: &Interface{Methods: []*Method{err}}, Obj: obj}
	}

	// predeclared constants
	for _, t := range predeclaredConstants {
		obj := def(go_ast.Con, t.name)
		obj.Type = Typ[t.kind]
		obj.Data = t.val
	}

	// predeclared functions
	for _, f := range predeclaredFunctions {
		def(go_ast.Fun, f.name).Type = f
	}

	universeIota = Universe.Lookup("iota")
}

// Objects with names containing blanks are internal and not entered into
// a scope. Objects with exported names are inserted in the unsafe package
// scope; other objects are inserted in the universe scope.
//
func def(kind go_ast.ObjKind, name string) *go_ast.Object {
	obj := go_ast.NewObj(kind, name)
	// insert non-internal objects into respective scope
	if strings.Index(name, " ") < 0 {
		scope := Universe
		// exported identifiers go into package unsafe
		if ast.IsExported(name) {
			scope = unsafe
		}
		if scope.Insert(obj) != nil {
			panic("internal error: double declaration")
		}
		obj.Decl = scope
	}
	return obj
}

// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// FILE UNDER CONSTRUCTION. ANY AND ALL PARTS MAY CHANGE.
// This file implements the universe and unsafe package scopes.

package compiler

import "go/ast"

var (
	scope    *ast.Scope // current scope to use for initialization
	Universe *ast.Scope
	Unsafe   *ast.Object // package unsafe
)

func define(kind ast.ObjKind, name string) *ast.Object {
	obj := ast.NewObj(kind, name)
	if scope.Insert(obj) != nil {
		panic("types internal error: double declaration")
	}
	obj.Decl = scope
	return obj
}

func defType(name string) {
	define(ast.Typ, name)

}

func defConst(name string) {
	obj := define(ast.Con, name)
	_ = obj // TODO(gri) fill in other properties
}

func defFun(name string) {
	obj := define(ast.Fun, name)
	_ = obj // TODO(gri) fill in other properties
}

func init() {
	scope = ast.NewScope(nil)
	Universe = scope

	defType("bool")
	defType("byte") // TODO(gri) should be an alias for uint8
	defType("rune") // TODO(gri) should be an alias for int
	defType("complex64")
	defType("complex128")
	defType("error")
	defType("float32")
	defType("float64")
	defType("int8")
	defType("int16")
	defType("int32")
	defType("int64")
	defType("string")
	defType("uint8")
	defType("uint16")
	defType("uint32")
	defType("uint64")
	defType("int")
	defType("uint")
	defType("uintptr")

	defConst("true")
	defConst("false")
	defConst("iota")
	defConst("nil")

	defFun("append")
	defFun("cap")
	defFun("close")
	defFun("complex")
	defFun("copy")
	defFun("delete")
	defFun("imag")
	defFun("len")
	defFun("make")
	defFun("new")
	defFun("panic")
	defFun("print")
	defFun("println")
	defFun("real")
	defFun("recover")

	scope = ast.NewScope(nil)
	Unsafe = ast.NewObj(ast.Pkg, "unsafe")
	Unsafe.Data = scope

	defType("Pointer")

	defFun("Alignof")
	defFun("Offsetof")
	defFun("Sizeof")
}

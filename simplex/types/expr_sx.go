package types

import (
	"github.com/fd/w/simplex/ast"
)

// rawExpr typechecks expression e and initializes x with the expression
// value or type. If an error occurred, x.mode is set to invalid.
// A hint != nil is used as operand type for untyped shifted operands;
// iota >= 0 indicates that the expression is part of a constant declaration.
// cycleOk indicates whether it is ok for a type expression to refer to itself.
//
// see expr.go:641
func (check *checker) rawExpr(x *operand, e ast.Expr, hint Type, iota int, cycleOk bool) {
	if trace {
		c := ""
		if cycleOk {
			c = " ⨁"
		}
		check.trace(e.Pos(), "%s (%s, %d%s)", e, typeString(hint), iota, c)
		defer check.untrace("=> %s", x)
	}

	if check.ctxt.Expr != nil {
		defer check.callExpr(x)
	}

	switch e := e.(type) {
	case *ast.BadExpr:
		goto Error // error was reported before

	case *ast.Ident:
		if e.Name == "_" {
			check.invalidOp(e.Pos(), "cannot use _ as value or type")
			goto Error
		}
		obj := e.Obj
		if obj == nil {
			goto Error // error was reported before
		}
		if obj.Type == nil {
			check.object(obj, cycleOk)
		}
		switch obj.Kind {
		case ast.Bad:
			goto Error // error was reported before
		case ast.Pkg:
			check.errorf(e.Pos(), "use of package %s not in selector", obj.Name)
			goto Error
		case ast.Con:
			if obj.Data == nil {
				goto Error // cycle detected
			}
			x.mode = constant
			if obj == universeIota {
				if iota < 0 {
					check.invalidAST(e.Pos(), "cannot use iota outside constant declaration")
					goto Error
				}
				x.val = int64(iota)
			} else {
				x.val = obj.Data
			}
		case ast.Typ:
			x.mode = typexpr
			if !cycleOk && underlying(obj.Type.(Type)) == nil {
				check.errorf(obj.Pos(), "illegal cycle in declaration of %s", obj.Name)
				x.expr = e
				x.typ = Typ[Invalid]
				return // don't goto Error - need x.mode == typexpr
			}
		case ast.Var:
			x.mode = variable
		case ast.Fun:
			x.mode = value
		default:
			unreachable()
		}
		x.typ = obj.Type.(Type)

	case *ast.Ellipsis:
		// ellipses are handled explicitly where they are legal
		// (array composite literals and parameter lists)
		check.errorf(e.Pos(), "invalid use of '...'")
		goto Error

	case *ast.BasicLit:
		x.setConst(e.Kind, e.Value)
		if x.mode == invalid {
			check.invalidAST(e.Pos(), "invalid literal %v", e.Value)
			goto Error
		}

	case *ast.FuncLit:
		if sig, ok := check.typ(e.Type, false).(*Signature); ok {
			x.mode = value
			x.typ = sig
			check.later(nil, sig, e.Body)
		} else {
			check.invalidAST(e.Pos(), "invalid function literal %s", e)
			goto Error
		}

	case *ast.CompositeLit:
		typ := hint
		openArray := false
		if e.Type != nil {
			// [...]T array types may only appear with composite literals.
			// Check for them here so we don't have to handle ... in general.
			typ = nil
			if atyp, _ := e.Type.(*ast.ArrayType); atyp != nil && atyp.Len != nil {
				if ellip, _ := atyp.Len.(*ast.Ellipsis); ellip != nil && ellip.Elt == nil {
					// We have an "open" [...]T array type.
					// Create a new ArrayType with unknown length (-1)
					// and finish setting it up after analyzing the literal.
					typ = &Array{Len: -1, Elt: check.typ(atyp.Elt, cycleOk)}
					openArray = true
				}
			}
			if typ == nil {
				typ = check.typ(e.Type, false)
			}
		}
		if typ == nil {
			check.errorf(e.Pos(), "missing type in composite literal")
			goto Error
		}

		switch utyp := underlying(deref(typ)).(type) {
		case *Struct:
			if len(e.Elts) == 0 {
				break
			}
			fields := utyp.Fields
			if _, ok := e.Elts[0].(*ast.KeyValueExpr); ok {
				// all elements must have keys
				visited := make([]bool, len(fields))
				for _, e := range e.Elts {
					kv, _ := e.(*ast.KeyValueExpr)
					if kv == nil {
						check.errorf(e.Pos(), "mixture of field:value and value elements in struct literal")
						continue
					}
					key, _ := kv.Key.(*ast.Ident)
					if key == nil {
						check.errorf(kv.Pos(), "invalid field name %s in struct literal", kv.Key)
						continue
					}
					i := utyp.fieldIndex(key.Name)
					if i < 0 {
						check.errorf(kv.Pos(), "unknown field %s in struct literal", key.Name)
						continue
					}
					// 0 <= i < len(fields)
					if visited[i] {
						check.errorf(kv.Pos(), "duplicate field name %s in struct literal", key.Name)
						continue
					}
					visited[i] = true
					check.expr(x, kv.Value, nil, iota)
					etyp := fields[i].Type
					if !x.isAssignable(etyp) {
						check.errorf(x.pos(), "cannot use %s as %s value in struct literal", x, etyp)
						continue
					}
				}
			} else {
				// no element must have a key
				for i, e := range e.Elts {
					if kv, _ := e.(*ast.KeyValueExpr); kv != nil {
						check.errorf(kv.Pos(), "mixture of field:value and value elements in struct literal")
						continue
					}
					check.expr(x, e, nil, iota)
					if i >= len(fields) {
						check.errorf(x.pos(), "too many values in struct literal")
						break // cannot continue
					}
					// i < len(fields)
					etyp := fields[i].Type
					if !x.isAssignable(etyp) {
						check.errorf(x.pos(), "cannot use %s as %s value in struct literal", x, etyp)
						continue
					}
				}
				if len(e.Elts) < len(fields) {
					check.errorf(e.Rbrace, "too few values in struct literal")
					// ok to continue
				}
			}

		case *Array:
			n := check.indexedElts(e.Elts, utyp.Elt, utyp.Len, iota)
			// if we have an "open" [...]T array, set the length now that we know it
			if openArray {
				utyp.Len = n
			}

		case *Slice:
			check.indexedElts(e.Elts, utyp.Elt, -1, iota)

		case *Map:
			visited := make(map[interface{}]bool, len(e.Elts))
			for _, e := range e.Elts {
				kv, _ := e.(*ast.KeyValueExpr)
				if kv == nil {
					check.errorf(e.Pos(), "missing key in map literal")
					continue
				}
				check.compositeLitKey(kv.Key)
				check.expr(x, kv.Key, nil, iota)
				if !x.isAssignable(utyp.Key) {
					check.errorf(x.pos(), "cannot use %s as %s key in map literal", x, utyp.Key)
					continue
				}
				if x.mode == constant {
					if visited[x.val] {
						check.errorf(x.pos(), "duplicate key %s in map literal", x.val)
						continue
					}
					visited[x.val] = true
				}
				check.expr(x, kv.Value, utyp.Elt, iota)
				if !x.isAssignable(utyp.Elt) {
					check.errorf(x.pos(), "cannot use %s as %s value in map literal", x, utyp.Elt)
					continue
				}
			}

		default:
			check.errorf(e.Pos(), "%s is not a valid composite literal type", typ)
			goto Error
		}

		x.mode = value
		x.typ = typ

	case *ast.ParenExpr:
		check.rawExpr(x, e.X, hint, iota, cycleOk)

	case *ast.SelectorExpr:
		sel := e.Sel.Name
		// If the identifier refers to a package, handle everything here
		// so we don't need a "package" mode for operands: package names
		// can only appear in qualified identifiers which are mapped to
		// selector expressions.
		if ident, ok := e.X.(*ast.Ident); ok {
			if obj := ident.Obj; obj != nil && obj.Kind == ast.Pkg {
				exp := obj.Data.(*ast.Scope).Lookup(sel)
				if exp == nil {
					check.errorf(e.Sel.Pos(), "cannot refer to unexported %s", sel)
					goto Error
				}
				// simplified version of the code for *ast.Idents:
				// imported objects are always fully initialized
				switch exp.Kind {
				case ast.Con:
					assert(exp.Data != nil)
					x.mode = constant
					x.val = exp.Data
				case ast.Typ:
					x.mode = typexpr
				case ast.Var:
					x.mode = variable
				case ast.Fun:
					x.mode = value
				default:
					unreachable()
				}
				x.expr = e
				x.typ = exp.Type.(Type)
				return
			}
		}

		check.exprOrType(x, e.X, nil, iota, false)
		if x.mode == invalid {
			goto Error
		}
		mode, typ := lookupField(x.typ, sel)
		if mode == invalid {
			check.invalidOp(e.Pos(), "%s has no single field or method %s", x, sel)
			goto Error
		}
		if x.mode == typexpr {
			// method expression
			sig, ok := typ.(*Signature)
			if !ok {
				check.invalidOp(e.Pos(), "%s has no method %s", x, sel)
				goto Error
			}
			// the receiver type becomes the type of the first function
			// argument of the method expression's function type
			// TODO(gri) at the moment, method sets don't correctly track
			// pointer vs non-pointer receivers => typechecker is too lenient
			x.mode = value
			x.typ = &Signature{
				Params:     append([]*Var{{"", x.typ}}, sig.Params...),
				Results:    sig.Results,
				IsVariadic: sig.IsVariadic,
			}
		} else {
			// regular selector
			x.mode = mode
			x.typ = typ
		}

	case *ast.IndexExpr:
		check.expr(x, e.X, hint, iota)

		valid := false
		length := int64(-1) // valid if >= 0
		switch typ := underlying(x.typ).(type) {
		case *Basic:
			if isString(typ) {
				valid = true
				if x.mode == constant {
					length = int64(len(x.val.(string)))
				}
				// an indexed string always yields a byte value
				// (not a constant) even if the string and the
				// index are constant
				x.mode = value
				x.typ = Typ[Byte]
			}

		case *Array:
			valid = true
			length = typ.Len
			if x.mode != variable {
				x.mode = value
			}
			x.typ = typ.Elt

		case *Pointer:
			if typ, _ := underlying(typ.Base).(*Array); typ != nil {
				valid = true
				length = typ.Len
				x.mode = variable
				x.typ = typ.Elt
			}

		case *Slice:
			valid = true
			x.mode = variable
			x.typ = typ.Elt

		case *Map:
			var key operand
			check.expr(&key, e.Index, nil, iota)
			if key.mode == invalid || !key.isAssignable(typ.Key) {
				check.invalidOp(x.pos(), "cannot use %s as map index of type %s", &key, typ.Key)
				goto Error
			}
			x.mode = valueok
			x.typ = typ.Elt
			x.expr = e
			return
		}

		if !valid {
			check.invalidOp(x.pos(), "cannot index %s", x)
			goto Error
		}

		if e.Index == nil {
			check.invalidAST(e.Pos(), "missing index expression for %s", x)
			return
		}

		check.index(e.Index, length, iota)
		// ok to continue

	case *ast.SliceExpr:
		check.expr(x, e.X, hint, iota)

		valid := false
		length := int64(-1) // valid if >= 0
		switch typ := underlying(x.typ).(type) {
		case *Basic:
			if isString(typ) {
				valid = true
				if x.mode == constant {
					length = int64(len(x.val.(string))) + 1 // +1 for slice
				}
				// a sliced string always yields a string value
				// of the same type as the original string (not
				// a constant) even if the string and the indices
				// are constant
				x.mode = value
				// x.typ doesn't change
			}

		case *Array:
			valid = true
			length = typ.Len + 1 // +1 for slice
			if x.mode != variable {
				check.invalidOp(x.pos(), "cannot slice %s (value not addressable)", x)
				goto Error
			}
			x.typ = &Slice{Elt: typ.Elt}

		case *Pointer:
			if typ, _ := underlying(typ.Base).(*Array); typ != nil {
				valid = true
				length = typ.Len + 1 // +1 for slice
				x.mode = variable
				x.typ = &Slice{Elt: typ.Elt}
			}

		case *Slice:
			valid = true
			x.mode = variable
			// x.typ doesn't change
		}

		if !valid {
			check.invalidOp(x.pos(), "cannot slice %s", x)
			goto Error
		}

		lo := int64(0)
		if e.Low != nil {
			lo = check.index(e.Low, length, iota)
		}

		hi := int64(-1)
		if e.High != nil {
			hi = check.index(e.High, length, iota)
		} else if length >= 0 {
			hi = length
		}

		if lo >= 0 && hi >= 0 && lo > hi {
			check.errorf(e.Low.Pos(), "inverted slice range: %d > %d", lo, hi)
			// ok to continue
		}

	case *ast.TypeAssertExpr:
		check.expr(x, e.X, hint, iota)
		if x.mode == invalid {
			goto Error
		}
		var T *Interface
		if T, _ = underlying(x.typ).(*Interface); T == nil {
			check.invalidOp(x.pos(), "%s is not an interface", x)
			goto Error
		}
		// x.(type) expressions are handled explicitly in type switches
		if e.Type == nil {
			check.errorf(e.Pos(), "use of .(type) outside type switch")
			goto Error
		}
		typ := check.typ(e.Type, false)
		if typ == Typ[Invalid] {
			goto Error
		}
		if method, wrongType := missingMethod(typ, T); method != nil {
			var msg string
			if wrongType {
				msg = "%s cannot have dynamic type %s (wrong type for method %s)"
			} else {
				msg = "%s cannot have dynamic type %s (missing method %s)"
			}
			check.errorf(e.Type.Pos(), msg, x, typ, method.Name)
			// ok to continue
		}
		x.mode = valueok
		x.expr = e
		x.typ = typ

	case *ast.CallExpr:
		check.exprOrType(x, e.Fun, nil, iota, false)
		if x.mode == invalid {
			goto Error
		} else if x.mode == typexpr {
			check.conversion(x, e, x.typ, iota)
		} else if sig, ok := underlying(x.typ).(*Signature); ok {
			// check parameters

			// If we have a trailing ... at the end of the parameter
			// list, the last argument must match the parameter type
			// []T of a variadic function parameter x ...T.
			passSlice := false
			if e.Ellipsis.IsValid() {
				if sig.IsVariadic {
					passSlice = true
				} else {
					check.errorf(e.Ellipsis, "cannot use ... in call to %s", e.Fun)
					// ok to continue
				}
			}

			// If we have a single argument that is a function call
			// we need to handle it separately. Determine if this
			// is the case without checking the argument.
			var call *ast.CallExpr
			if len(e.Args) == 1 {
				call, _ = unparen(e.Args[0]).(*ast.CallExpr)
			}

			n := 0 // parameter count
			if call != nil {
				// We have a single argument that is a function call.
				check.expr(x, call, nil, -1)
				if x.mode == invalid {
					goto Error // TODO(gri): we can do better
				}
				if t, _ := x.typ.(*Result); t != nil {
					// multiple result values
					n = len(t.Values)
					for i, obj := range t.Values {
						x.mode = value
						x.expr = nil // TODO(gri) can we do better here? (for good error messages)
						x.typ = obj.Type.(Type)
						check.argument(sig, i, nil, x, passSlice && i+1 == n)
					}
				} else {
					// single result value
					n = 1
					check.argument(sig, 0, nil, x, passSlice)
				}

			} else {
				// We don't have a single argument or it is not a function call.
				n = len(e.Args)
				for i, arg := range e.Args {
					check.argument(sig, i, arg, x, passSlice && i+1 == n)
				}
			}

			// determine if we have enough arguments
			if sig.IsVariadic {
				// a variadic function accepts an "empty"
				// last argument: count one extra
				n++
			}
			if n < len(sig.Params) {
				check.errorf(e.Fun.Pos(), "too few arguments in call to %s", e.Fun)
				// ok to continue
			}

			// determine result
			switch len(sig.Results) {
			case 0:
				x.mode = novalue
			case 1:
				x.mode = value
				x.typ = sig.Results[0].Type.(Type)
			default:
				x.mode = value
				x.typ = &Result{Values: sig.Results}
			}

		} else if bin, ok := x.typ.(*builtin); ok {
			check.builtin(x, e, bin, iota)

		} else {
			check.invalidOp(x.pos(), "cannot call non-function %s", x)
			goto Error
		}

	case *ast.StarExpr:
		check.exprOrType(x, e.X, hint, iota, true)
		switch x.mode {
		case invalid:
			goto Error
		case typexpr:
			x.typ = &Pointer{Base: x.typ}
		default:
			if typ, ok := x.typ.(*Pointer); ok {
				x.mode = variable
				x.typ = typ.Base
			} else {
				check.invalidOp(x.pos(), "cannot indirect %s", x)
				goto Error
			}
		}

	case *ast.UnaryExpr:
		check.expr(x, e.X, hint, iota)
		check.unary(x, e.Op)

	case *ast.BinaryExpr:
		var y operand
		check.expr(x, e.X, hint, iota)
		check.expr(&y, e.Y, hint, iota)
		check.binary(x, &y, e.Op, hint)

	case *ast.KeyValueExpr:
		// key:value expressions are handled in composite literals
		check.invalidAST(e.Pos(), "no key:value expected")
		goto Error

	case *ast.ArrayType:
		if e.Len != nil {
			check.expr(x, e.Len, nil, iota)
			if x.mode == invalid {
				goto Error
			}
			if x.mode != constant {
				if x.mode != invalid {
					check.errorf(x.pos(), "array length %s must be constant", x)
				}
				goto Error
			}
			n, ok := x.val.(int64)
			if !ok || n < 0 {
				check.errorf(x.pos(), "invalid array length %s", x)
				goto Error
			}
			x.typ = &Array{Len: n, Elt: check.typ(e.Elt, cycleOk)}
		} else {
			x.typ = &Slice{Elt: check.typ(e.Elt, true)}
		}
		x.mode = typexpr

	case *ast.StructType:
		x.mode = typexpr
		x.typ = &Struct{Fields: check.collectFields(e.Fields, cycleOk)}

	case *ast.FuncType:
		params, isVariadic := check.collectParams(e.Params, true)
		results, _ := check.collectParams(e.Results, false)
		x.mode = typexpr
		x.typ = &Signature{Recv: nil, Params: params, Results: results, IsVariadic: isVariadic}

	case *ast.InterfaceType:
		x.mode = typexpr
		x.typ = &Interface{Methods: check.collectMethods(e.Methods)}

	case *ast.MapType:
		x.mode = typexpr
		x.typ = &Map{Key: check.typ(e.Key, true), Elt: check.typ(e.Value, true)}

	case *ast.ChanType:
		x.mode = typexpr
		x.typ = &Chan{Dir: e.Dir, Elt: check.typ(e.Value, true)}

	//=== start custom
	case *ast.StepExpr:
		check.sx_step(x, e, hint, iota)
		if x.mode == invalid {
			goto Error
		}

	case *ast.ViewType:
		x.mode = typexpr
		if e.Key != nil {
			x.typ = &View{Key: check.typ(e.Key, false), Elt: check.typ(e.Value, false)}
		} else {
			x.typ = &View{Elt: check.typ(e.Value, false)}
		}

	case *ast.TableType:
		x.mode = typexpr
		x.typ = &Table{Key: check.typ(e.Key, false), Elt: check.typ(e.Value, false)}
	//=== end custom

	default:
		check.dump("e = %s", e)
		unreachable()
	}

	// everything went well
	x.expr = e
	return

Error:
	x.mode = invalid
	x.expr = e
}

func (check *checker) sx_step(x *operand, e *ast.StepExpr, hint Type, iota int) {
	switch e.StepType {

	case ast.SelectStep:
		check.sx_predicate_function(e, e.X, e.F, x, hint, iota, "select")

	case ast.RejectStep:
		check.sx_predicate_function(e, e.X, e.F, x, hint, iota, "reject")

	case ast.DetectStep:
		check.sx_predicate_function(e, e.X, e.F, x, hint, iota, "detect")
		if x.mode != invalid {
			x.typ = x.typ.(*View).Elt
		}

	case ast.CollectStep:
		res_typ := check.sx_map_function(e, e.X, e.F, x, hint, iota, "collect")
		if x.mode != invalid {
			view_typ := x.typ.(*View)
			x.typ = &View{Key: view_typ.Key, Elt: res_typ}
		}

	case ast.InjectStep:
		res_typ := check.sx_inject_function(e, e.X, e.F, x, hint, iota, "inject")
		if x.mode != invalid {
			view_typ := x.typ.(*View)
			x.typ = &View{Key: view_typ.Key, Elt: res_typ}
		}

	case ast.GroupStep:
		res_typ := check.sx_map_function(e, e.X, e.F, x, hint, iota, "group")
		if x.mode != invalid {
			view_typ := x.typ.(*View)
			x.typ = &View{Key: res_typ, Elt: view_typ}
		}

	case ast.IndexStep:
		res_typ := check.sx_map_function(e, e.X, e.F, x, hint, iota, "index")
		if x.mode != invalid {
			view_typ := x.typ.(*View)
			x.typ = &View{Key: res_typ, Elt: view_typ.Elt}
		}

	case ast.SortStep:
		check.sx_map_function(e, e.X, e.F, x, hint, iota, "sort")
		if x.mode != invalid {
			view_typ := x.typ.(*View)
			x.typ = &View{Elt: view_typ.Elt}
		}

	default:
		unreachable()

	}
}

func (check *checker) sx_predicate_function(e, x, f ast.Expr, op *operand, hint Type, iota int, expr string) {
	check.expr(op, x, nil, iota)
	view_typ, ok := op.typ.(*View)
	if !ok {
		check.errorf(x.Pos(), "not a view receiver %s", x)
		op.mode = invalid
		return
	}

	var f_op operand
	check.expr(&f_op, f, nil, iota)
	if sig, ok := f_op.typ.(*Signature); ok {
		if len(sig.Results) != 1 {
			check.errorf(f.Pos(), "not a "+expr+" function %s", f)
			op.mode = invalid
			return
		}
		if !isBoolean(sig.Results[0].Type) {
			check.errorf(f.Pos(), "not a "+expr+" function %s", f)
			op.mode = invalid
			return
		}
		if len(sig.Params) != 1 {
			check.errorf(f.Pos(), "not a "+expr+" function %s", f)
			op.mode = invalid
			return
		}
		if !isIdentical(sig.Params[0].Type, view_typ.Elt) {
			check.errorf(f.Pos(), "not a "+expr+" function %s", f)
			op.mode = invalid
			return
		}
	} else {
		check.errorf(f.Pos(), "not a "+expr+" function %s", f)
		op.mode = invalid
		return
	}
}

func (check *checker) sx_map_function(e, x, f ast.Expr, op *operand, hint Type, iota int, expr string) Type {
	var res_typ Type

	check.expr(op, x, nil, iota)
	view_typ, ok := op.typ.(*View)

	if !ok {
		check.errorf(x.Pos(), "not a view receiver %s", x)
		op.mode = invalid
		return nil
	}

	var f_op operand
	check.expr(&f_op, f, nil, iota)
	if sig, ok := f_op.typ.(*Signature); ok {
		if len(sig.Results) != 1 {
			check.errorf(f.Pos(), "not a "+expr+" function %s", f)
			op.mode = invalid
			return nil
		}
		if _, isSig := sig.Results[0].Type.(*Signature); isSig {
			check.errorf(f.Pos(), "not a "+expr+" function %s", f)
			op.mode = invalid
			return nil
		}
		res_typ = sig.Results[0].Type

		if len(sig.Params) != 1 {
			check.errorf(f.Pos(), "not a "+expr+" function %s", f)
			op.mode = invalid
			return nil
		}
		if !isIdentical(sig.Params[0].Type, view_typ.Elt) {
			check.errorf(f.Pos(), "not a "+expr+" function %s", f)
			op.mode = invalid
			return nil
		}

	} else {
		check.errorf(f.Pos(), "not a "+expr+" function %s", f)
		op.mode = invalid
		return nil
	}

	return res_typ
}

func (check *checker) sx_inject_function(e, x, f ast.Expr, op *operand, hint Type, iota int, expr string) Type {
	var res_typ Type

	check.expr(op, x, hint, iota)
	view_typ, ok := op.typ.(*View)

	if !ok {
		check.errorf(x.Pos(), "not a view receiver %s", x)
		op.mode = invalid
		return nil
	}

	group_view_typ, ok := view_typ.Elt.(*View)
	if !ok {
		check.errorf(x.Pos(), "not an inject view receiver %s", x)
		op.mode = invalid
		return nil
	}

	var f_op operand
	check.expr(&f_op, f, hint, iota)
	if sig, ok := f_op.typ.(*Signature); ok {
		if len(sig.Results) != 1 {
			check.errorf(f.Pos(), "not a "+expr+" function %s", f)
			op.mode = invalid
			return nil
		}
		if _, isSig := sig.Results[0].Type.(*Signature); isSig {
			check.errorf(f.Pos(), "not a "+expr+" function %s", f)
			op.mode = invalid
			return nil
		}
		res_typ = sig.Results[0].Type

		if len(sig.Params) != 2 {
			check.errorf(f.Pos(), "not a "+expr+" function %s", f)
			op.mode = invalid
			return nil
		}

		acc_typ := &Slice{Elt: res_typ}
		if !isIdentical(sig.Params[0].Type, acc_typ) {
			check.errorf(f.Pos(), "not a "+expr+" function %s", f)
			op.mode = invalid
			return nil
		}
		if !isIdentical(sig.Params[1].Type, group_view_typ.Elt) {
			check.errorf(f.Pos(), "not a "+expr+" function %s", f)
			op.mode = invalid
			return nil
		}

	} else {
		check.errorf(f.Pos(), "not a "+expr+" function %s", f)
		op.mode = invalid
		return nil
	}

	return res_typ
}

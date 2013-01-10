package compiler

import (
	"go/ast"
	//"fmt"
)

func resolve_view_type(scope *ast.Scope, node interface{}) *ViewDecl {
	switch n := node.(type) {

	case *ast.CallExpr:

		switch fun := n.Fun.(type) {

		//case *ast.Ident:
		//fun.Obj

		case *ast.SelectorExpr:
			if fun.Sel.Obj == nil {
				return nil
			}

			func_decl, ok := fun.Sel.Obj.Decl.(*ast.FuncDecl)
			if !ok {
				return nil
			}

			results := func_decl.Type.Results
			if results == nil {
				return nil
			}

			if len(results.List) == 0 {
				return nil
			}

			field := results.List[0]
			if len(field.Names) > 1 {
				return nil
			}

			typ := field.Type
			return resolve_view_type(scope, typ)

		default:
			ast.Print(nil, fun)
		}

	case *ast.Ident:
		if n.Obj == nil {
			return nil
		}

		if view_decl, ok := n.Obj.Data.(*ViewDecl); ok {
			return view_decl
		}

		switch n.Obj.Kind {

		case ast.Var:
			// resolve var type
			var expr ast.Expr
			switch decl := n.Obj.Decl.(type) {

			case *ast.ValueSpec:
				for i, name := range decl.Names {
					if name.Name == n.Name {
						expr = decl.Values[i]
						break
					}
				}

			case *ast.Field:
				expr = decl.Type

			default:
				panic("unsupported node")

			}

			return resolve_view_type(scope, expr)

		case ast.Fun:
			decl := n.Obj.Decl.(*ast.FuncDecl)
			results := decl.Type.Results
			if results == nil {
				return nil
			}
			if len(results.List) != 1 {
				return nil
			}
			field := results.List[0]
			if len(field.Names) > 1 {
				return nil
			}
			return resolve_view_type(scope, field.Type)

		default:
			ast.Print(nil, n)
		}

		// T.collect(F(t)t') => T'
		// T.inject(F(t, a)a) => A
		// T.select(...) => T
		// T.reject(...) => T
		// T.detect(...) => T
		// T.sort(...) => T
		// T.group(F(t)g) => G{g, T}
		// x.F(...)T => T

	default:
		ast.Print(nil, n)
		return nil
	}

	ast.Print(nil, node)
	return nil
}

func resolve_function_type(scope *ast.Scope, node interface{}) *ViewDecl {
	switch n := node.(type) {

	case *ast.Ident:
		if n.Obj == nil {
			return nil
		}

		switch n.Obj.Kind {

		case ast.Fun:
			decl := n.Obj.Decl.(*ast.FuncDecl)
			results := decl.Type.Results
			if results == nil {
				return nil
			}
			if len(results.List) != 1 {
				return nil
			}
			field := results.List[0]
			if len(field.Names) > 1 {
				return nil
			}
			return resolve_view_type(scope, field.Type)

		default:
			ast.Print(nil, n)
		}

	//case *ast.SelectorExpr:
	// T.collect(F(t)t') => T'
	// T.inject(F(t, a)a) => A
	// T.select(...) => T
	// T.reject(...) => T
	// T.detect(...) => T
	// T.sort(...) => T
	// T.group(F(t)g) => G{g, T}
	// x.F(...)T => T

	default:
		ast.Print(nil, n)
	}

	return nil
}

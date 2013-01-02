package compiler

import (
	"go/ast"
	//"fmt"
)

func resolve_type(scope *ast.Scope, node ast.Node) (o *ast.Object) {
	defer func() {
		if o != nil {
			if _, k := o.Data.(*ViewDecl); k && o.Type == nil {
				o.Type = o.Data
			}
		}
	}()

	switch n := node.(type) {

	case *ast.CallExpr:
		obj := resolve_type(scope, n.Fun)
		return obj

	case *ast.Field:
		return resolve_type(scope, n.Type)

	case *ast.Ident:
		if n.Obj == nil {
			if !resolve(scope, n) {
				return nil
			}
		}
		if n.Obj.Type != nil {
			return n.Obj
		}
		if node, ok := n.Obj.Decl.(ast.Node); ok {
			typ := resolve_type(scope, node)
			if typ != nil {
				n.Obj.Type = typ.Type
			}
			return n.Obj
		}
		if view, ok := n.Obj.Data.(*ViewDecl); ok {
			n.Obj.Type = view
			return n.Obj
		}
		return nil

	case *ast.SelectorExpr:
		if n.Sel.Obj != nil {
			if _, ok := n.Sel.Obj.Data.(*ViewDecl); ok {
				return n.Sel.Obj
			}
		}

		x_obj := resolve_type(scope, n.X)
		if x_obj == nil {
			return nil
		}

		switch x_obj.Kind {

		case ast.Pkg:
			return resolve_type(x_obj.Data.(*ast.Scope), n.Sel)

		case ast.Fun:
			decl := x_obj.Decl.(*ast.FuncDecl)
			r := decl.Type.Results
			if r == nil {
				return nil
			}
			if r.NumFields() != 1 {
				return nil
			}
			typ_spec := r.List[0].Type
			obj := resolve_type(scope, typ_spec)
			return obj

		case ast.Var:
			decl := x_obj.Decl.(*ast.Field)
			typ_spec := decl.Type
			obj := resolve_type(scope, typ_spec)
			return obj

		case ast.Typ:
			return x_obj

		default:
			ast.Print(nil, node)

		}

	case *ast.ValueSpec:
		for i, val := range n.Names {
			obj := resolve_type(scope, n.Values[i])
			val.Obj.Type = obj.Type
		}
		return nil

	default:
		//ast.Print(nil, node)

	}

	return nil
}

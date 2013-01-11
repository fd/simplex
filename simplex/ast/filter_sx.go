package ast

// see filter.go:127
func filterType(typ Expr, f Filter, export bool) bool {
	switch t := typ.(type) {
	case *Ident:
		return f(t.Name)
	case *ParenExpr:
		return filterType(t.X, f, export)
	case *ArrayType:
		return filterType(t.Elt, f, export)
	case *StructType:
		if filterFieldList(t.Fields, f, export) {
			t.Incomplete = true
		}
		return len(t.Fields.List) > 0
	case *FuncType:
		b1 := filterParamList(t.Params, f, export)
		b2 := filterParamList(t.Results, f, export)
		return b1 || b2
	case *InterfaceType:
		if filterFieldList(t.Methods, f, export) {
			t.Incomplete = true
		}
		return len(t.Methods.List) > 0
	case *MapType:
		b1 := filterType(t.Key, f, export)
		b2 := filterType(t.Value, f, export)
		return b1 || b2
	case *ChanType:
		return filterType(t.Value, f, export)

	//=== start custom
	case *ViewType:
		if t.Key == nil {
			return filterType(t.Value, f, export)
		}
		b1 := filterType(t.Key, f, export)
		b2 := filterType(t.Value, f, export)
		return b1 || b2
	case *TableType:
		b1 := filterType(t.Key, f, export)
		b2 := filterType(t.Value, f, export)
		return b1 || b2
		//=== end custom

	}
	return false
}

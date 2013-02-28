package compiler

import (
	"simplex.sh/lang/ast"
	"simplex.sh/lang/token"
	"simplex.sh/lang/types"
	"strconv"
)

func (c *Context) convert_sx_to_go() error {

	for _, name := range c.SxFiles {
		var (
			file         = c.AstFiles[name]
			runtime_name = ""
			cas_name     = ""
			cas_missing  bool
		)

		for _, imp := range file.Imports {
			if imp.Path.Value == `"simplex.sh/runtime"` {
				if imp.Name == nil {
					runtime_name = "runtime"
				} else {
					runtime_name = imp.Name.Name
				}
			}
			if imp.Path.Value == `"simplex.sh/cas"` {
				if imp.Name == nil {
					cas_name = "cas"
				} else {
					cas_name = imp.Name.Name
				}
			}
		}
		if runtime_name == "" {
			runtime_name = "sx_runtime"
			add_import(file, runtime_name, "simplex.sh/runtime")
		}
		if cas_name == "" {
			cas_name = "sx_cas"
			cas_missing = true
		}

		v := &builtin_function_conv{
			c.NodeTypes,
			c.FileSet,
			runtime_name,
			cas_name,
			c.ImportPath,
			false,
		}
		ast.Replace(v, file)

		if cas_missing && v.used_cas_import {
			add_import(file, cas_name, "simplex.sh/cas")
		}
	}

	return nil
}

func add_import(file *ast.File, name, path string) {
	imp := &ast.ImportSpec{
		Name: ast.NewIdent(name),
		Path: &ast.BasicLit{
			Kind:  token.STRING,
			Value: strconv.Quote(path),
		},
	}
	decl := &ast.GenDecl{
		Tok:   token.IMPORT,
		Specs: []ast.Spec{imp},
	}
	file.Decls = append([]ast.Decl{decl}, file.Decls...)
	file.Imports = append(file.Imports, imp)
}

type (
	builtin_function_conv struct {
		mapping         map[ast.Node]types.Type
		fset            *token.FileSet
		runtime_name    string
		cas_name        string
		package_name    string
		used_cas_import bool
	}
)

func (conv *builtin_function_conv) Replace(node ast.Node) (ast.Replacer, ast.Node) {
	switch n := node.(type) {

	case *ast.ArrayType:
		n.Elt = conv.convert_type(n.Elt)

		//case *ast.AssignStmt:
		//case *ast.BadDecl:
		//case *ast.BadExpr:
		//case *ast.BadStmt:
		//case *ast.BasicLit:
		//case *ast.BinaryExpr:
		//case *ast.BlockStmt:
		//case *ast.BranchStmt:

	case *ast.CallExpr:
		if ident, ok := n.Fun.(*ast.Ident); ok {

			switch ident.Name {
			case "make":
				conv.convert_make(n)

			default:
				if ident.Obj != nil && ident.Obj.Kind == ast.Typ {
					typ := conv.mapping[ident]
					if n_typ, ok := typ.(*types.NamedType); ok {
						if _, ok := n_typ.Underlying.(types.Viewish); ok {
							r := &ast.CompositeLit{
								Type: ident,
								Elts: n.Args,
							}
							return conv, r
						}
					}
				}

			}
		}

		if sel, ok := n.Fun.(*ast.SelectorExpr); ok {
			switch sel.Sel.Name {

			case "select":
				conv.convert_select(n)

			case "reject":
				conv.convert_reject(n)

			case "collect":
				conv.convert_collect(n)

			case "inject":
				conv.convert_inject(n)

			case "group":
				conv.convert_group(n)

			case "index":
				conv.convert_index(n)

			case "sort":
				conv.convert_sort(n)

			}
		}

	//case *ast.CaseClause:

	case *ast.ChanType:
		n.Value = conv.convert_type(n.Value)

	//case *ast.CommClause:
	//case *ast.Comment:
	//case *ast.CommentGroup:

	case *ast.CompositeLit:
		n.Type = conv.convert_type(n.Type)

	//case *ast.DeclStmt:
	//case *ast.DeferStmt:

	case *ast.Ellipsis:
		n.Elt = conv.convert_type(n.Elt)

	//case *ast.EmptyStmt:
	//case *ast.ExprStmt:

	case *ast.Field:
		n.Type = conv.convert_type(n.Type)

	//case *ast.FieldList:
	//case *ast.File:
	//case *ast.ForStmt:
	//case *ast.FuncDecl:
	//case *ast.FuncLit:
	//case *ast.FuncType:
	//case *ast.GenDecl:
	//case *ast.GoStmt:
	//case *ast.Ident:
	//case *ast.IfStmt:
	//case *ast.ImportSpec:
	//case *ast.IncDecStmt:
	//case *ast.IndexExpr:
	//case *ast.InterfaceType:
	//case *ast.KeyValueExpr:
	//case *ast.LabeledStmt:

	case *ast.MapType:
		n.Key = conv.convert_type(n.Key)
		n.Value = conv.convert_type(n.Value)

	//case *ast.Package:
	//case *ast.ParenExpr:
	//case *ast.RangeStmt:
	//case *ast.ReturnStmt:
	//case *ast.SelectStmt:
	//case *ast.SelectorExpr:
	//case *ast.SendStmt:
	//case *ast.SliceExpr:
	//case *ast.StarExpr:
	//case *ast.StructType:
	//case *ast.SwitchStmt:

	case *ast.TableType:
		n.Key = conv.convert_type(n.Key)
		n.Value = conv.convert_type(n.Value)

	case *ast.TypeAssertExpr:
		n.Type = conv.convert_type(n.Type)

	case *ast.TypeSpec:
		switch n.Type.(type) {
		case *ast.ViewType, *ast.TableType:
			n.Type = &ast.StructType{
				Fields: &ast.FieldList{
					List: []*ast.Field{
						{Type: conv.convert_type(n.Type)},
					},
				},
			}
		default:
			n.Type = conv.convert_type(n.Type)
		}

	//case *ast.TypeSwitchStmt:
	//case *ast.UnaryExpr:

	case *ast.ValueSpec:
		n.Type = conv.convert_type(n.Type)

	case *ast.ViewType:
		n.Key = conv.convert_type(n.Key)
		n.Value = conv.convert_type(n.Value)

	}
	return conv, nil
}

func (conv *builtin_function_conv) convert_type(expr ast.Expr) ast.Expr {
	switch expr.(type) {
	case *ast.ViewType, *ast.TableType:
		typ := conv.mapping[expr]
		name := view_type_name(typ)
		return ast.NewIdent(name)
	}
	return expr
}

func (conv *builtin_function_conv) convert_make(call *ast.CallExpr) {
	typ := conv.mapping[call]
	if _, ok := typ.(*types.Table); !ok {
		return
	}

	call.Fun = ast.NewIdent("new_" + view_type_name(typ))
	call.Args = []ast.Expr{
		&ast.SelectorExpr{
			X:   ast.NewIdent(conv.runtime_name),
			Sel: ast.NewIdent("Env"),
		},
		&ast.BasicLit{
			Kind:  token.STRING,
			Value: strconv.QuoteToASCII(strconv.Quote(conv.package_name) + "." + sx_type_string(typ)),
		},
	}
}

func (conv *builtin_function_conv) convert_method(call *ast.CallExpr, name string) {
	var (
		recv  = call.Fun.(*ast.SelectorExpr).X
		args  = call.Args
		o_typ = conv.mapping[call]
		pos   = conv.fset.Position(call.Pos()).String()
	)

	args = append([]ast.Expr{recv}, args...)
	args = append(
		args,
		&ast.BasicLit{
			Kind: token.STRING,
			Value: strconv.QuoteToASCII(
				strconv.Quote(conv.package_name) + "." + sx_type_string(o_typ) + "[" +
					pos + "]",
			),
		},
	)

	call.Fun = ast.NewIdent("wrap_" + view_type_name(o_typ))
	call.Args = []ast.Expr{
		&ast.CallExpr{
			Fun: &ast.SelectorExpr{
				ast.NewIdent(conv.runtime_name),
				ast.NewIdent(name),
			},
			Args: args,
		},
	}
}

func (conv *builtin_function_conv) wrap_predicate_function(e ast.Expr) ast.Expr {
	typ := conv.mapping[e]
	if typ == nil {
		return nil
	}

	sig, ok := typ.(*types.Signature)
	if !ok {
		return nil
	}

	if sig.Recv != nil {
		return nil
	}

	if len(sig.Params) != 1 {
		return nil
	}

	if len(sig.Results) != 1 {
		return nil
	}

	conv.used_cas_import = true

	wrapper := &ast.FuncLit{
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							ast.NewIdent("sx_ctx"),
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent(conv.runtime_name),
								Sel: ast.NewIdent("Context"),
							},
						},
					},
					{
						Names: []*ast.Ident{
							ast.NewIdent("sx_m_addr"),
						},
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent(conv.cas_name),
							Sel: ast.NewIdent("Addr"),
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: ast.NewIdent("bool"),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.DeclStmt{
					Decl: &ast.GenDecl{
						Tok: token.VAR,
						Specs: []ast.Spec{
							&ast.ValueSpec{
								Names: []*ast.Ident{ast.NewIdent("sx_m")},
								Type:  ast.NewIdent(view_type_name(sig.Params[0].Type)),
							},
						},
					},
				},
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("sx_ctx"),
							Sel: ast.NewIdent("Load"),
						},
						Args: []ast.Expr{
							ast.NewIdent("sx_m_addr"),
							&ast.UnaryExpr{
								Op: token.AND,
								X:  ast.NewIdent("sx_m"),
							},
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.CallExpr{
							Fun: e,
							Args: []ast.Expr{
								ast.NewIdent("sx_m"),
							},
						},
					},
				},
			},
		},
	}

	return wrapper
}

func (conv *builtin_function_conv) wrap_map_function(e ast.Expr) ast.Expr {
	typ := conv.mapping[e]
	if typ == nil {
		return nil
	}

	sig, ok := typ.(*types.Signature)
	if !ok {
		return nil
	}

	if sig.Recv != nil {
		return nil
	}

	if len(sig.Params) != 1 {
		return nil
	}

	if len(sig.Results) != 1 {
		return nil
	}

	conv.used_cas_import = true

	wrapper := &ast.FuncLit{
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							ast.NewIdent("sx_ctx"),
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent(conv.runtime_name),
								Sel: ast.NewIdent("Context"),
							},
						},
					},
					{
						Names: []*ast.Ident{
							ast.NewIdent("sx_m_addr"),
						},
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent(conv.cas_name),
							Sel: ast.NewIdent("Addr"),
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent(conv.cas_name),
							Sel: ast.NewIdent("Addr"),
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.DeclStmt{
					Decl: &ast.GenDecl{
						Tok: token.VAR,
						Specs: []ast.Spec{
							&ast.ValueSpec{
								Names: []*ast.Ident{ast.NewIdent("sx_m")},
								Type:  ast.NewIdent(view_type_name(sig.Params[0].Type)),
							},
						},
					},
				},
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("sx_ctx"),
							Sel: ast.NewIdent("Load"),
						},
						Args: []ast.Expr{
							ast.NewIdent("sx_m_addr"),
							&ast.UnaryExpr{
								Op: token.AND,
								X:  ast.NewIdent("sx_m"),
							},
						},
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("sx_n"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: e,
							Args: []ast.Expr{
								ast.NewIdent("sx_m"),
							},
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("sx_ctx"),
								Sel: ast.NewIdent("Save"),
							},
							Args: []ast.Expr{
								&ast.UnaryExpr{
									Op: token.AND,
									X:  ast.NewIdent("sx_n"),
								},
							},
						},
					},
				},
			},
		},
	}

	return wrapper
}

func (conv *builtin_function_conv) wrap_sort_function(e ast.Expr) ast.Expr {
	typ := conv.mapping[e]
	if typ == nil {
		return nil
	}

	sig, ok := typ.(*types.Signature)
	if !ok {
		return nil
	}

	if sig.Recv != nil {
		return nil
	}

	if len(sig.Params) != 1 {
		return nil
	}

	if len(sig.Results) != 1 {
		return nil
	}

	conv.used_cas_import = true

	wrapper := &ast.FuncLit{
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							ast.NewIdent("sx_ctx"),
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent(conv.runtime_name),
								Sel: ast.NewIdent("Context"),
							},
						},
					},
					{
						Names: []*ast.Ident{
							ast.NewIdent("sx_m_addr"),
						},
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent(conv.cas_name),
							Sel: ast.NewIdent("Addr"),
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.InterfaceType{Methods: &ast.FieldList{}},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.DeclStmt{
					Decl: &ast.GenDecl{
						Tok: token.VAR,
						Specs: []ast.Spec{
							&ast.ValueSpec{
								Names: []*ast.Ident{ast.NewIdent("sx_m")},
								Type:  ast.NewIdent(view_type_name(sig.Params[0].Type)),
							},
						},
					},
				},
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("sx_ctx"),
							Sel: ast.NewIdent("Load"),
						},
						Args: []ast.Expr{
							ast.NewIdent("sx_m_addr"),
							&ast.UnaryExpr{
								Op: token.AND,
								X:  ast.NewIdent("sx_m"),
							},
						},
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("sx_n"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: e,
							Args: []ast.Expr{
								ast.NewIdent("sx_m"),
							},
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						ast.NewIdent("sx_n"),
					},
				},
			},
		},
	}

	return wrapper
}

func (conv *builtin_function_conv) convert_select(call *ast.CallExpr) {
	var (
		recv  = call.Fun.(*ast.SelectorExpr).X
		i_typ = conv.mapping[recv]
	)

	if _, ok := underlying_type(i_typ).(types.Viewish); !ok {
		return
	}

	if len(call.Args) != 1 {
		return
	}

	arg0 := conv.wrap_predicate_function(call.Args[0])
	if arg0 == nil {
		return
	}
	call.Args[0] = arg0

	conv.convert_method(call, "Select")
}

func (conv *builtin_function_conv) convert_reject(call *ast.CallExpr) {
	var (
		recv  = call.Fun.(*ast.SelectorExpr).X
		i_typ = conv.mapping[recv]
	)

	if _, ok := underlying_type(i_typ).(types.Viewish); !ok {
		return
	}

	if len(call.Args) != 1 {
		return
	}

	arg0 := conv.wrap_predicate_function(call.Args[0])
	if arg0 == nil {
		return
	}
	call.Args[0] = arg0

	conv.convert_method(call, "Reject")
}

func (conv *builtin_function_conv) convert_detect(call *ast.CallExpr) {
	var (
		recv  = call.Fun.(*ast.SelectorExpr).X
		i_typ = conv.mapping[recv]
	)

	if _, ok := underlying_type(i_typ).(types.Viewish); !ok {
		return
	}

	if len(call.Args) != 1 {
		return
	}

	arg0 := conv.wrap_predicate_function(call.Args[0])
	if arg0 == nil {
		return
	}
	call.Args[0] = arg0

	conv.convert_method(call, "Detect")
}

func (conv *builtin_function_conv) convert_collect(call *ast.CallExpr) {
	var (
		recv  = call.Fun.(*ast.SelectorExpr).X
		i_typ = conv.mapping[recv]
	)

	if _, ok := underlying_type(i_typ).(types.Viewish); !ok {
		return
	}

	if len(call.Args) != 1 {
		return
	}

	arg0 := conv.wrap_map_function(call.Args[0])
	if arg0 == nil {
		return
	}
	call.Args[0] = arg0

	conv.convert_method(call, "Collect")
}

func (conv *builtin_function_conv) convert_inject(call *ast.CallExpr) {
	var (
		recv  = call.Fun.(*ast.SelectorExpr).X
		i_typ = conv.mapping[recv]
	)

	switch t := i_typ.(type) {
	case *types.View:
	case *types.Table:
	case *types.NamedType:
		if _, ok := t.Underlying.(types.Viewish); !ok {
			return
		}
	default:
		return
	}

	conv.convert_method(call, "Inject")
}

func (conv *builtin_function_conv) convert_group(call *ast.CallExpr) {
	var (
		recv  = call.Fun.(*ast.SelectorExpr).X
		i_typ = conv.mapping[recv]
	)

	if _, ok := underlying_type(i_typ).(types.Viewish); !ok {
		return
	}

	if len(call.Args) != 1 {
		return
	}

	arg0 := conv.wrap_map_function(call.Args[0])
	if arg0 == nil {
		return
	}
	call.Args[0] = arg0

	conv.convert_method(call, "Group")
}

func (conv *builtin_function_conv) convert_index(call *ast.CallExpr) {
	var (
		recv  = call.Fun.(*ast.SelectorExpr).X
		i_typ = conv.mapping[recv]
	)

	if _, ok := underlying_type(i_typ).(types.Viewish); !ok {
		return
	}

	if len(call.Args) != 1 {
		return
	}

	arg0 := conv.wrap_map_function(call.Args[0])
	if arg0 == nil {
		return
	}
	call.Args[0] = arg0

	conv.convert_method(call, "Index")
}

func (conv *builtin_function_conv) convert_sort(call *ast.CallExpr) {
	var (
		recv  = call.Fun.(*ast.SelectorExpr).X
		i_typ = conv.mapping[recv]
	)

	if _, ok := underlying_type(i_typ).(types.Viewish); !ok {
		return
	}

	if len(call.Args) != 1 {
		return
	}

	arg0 := conv.wrap_sort_function(call.Args[0])
	if arg0 == nil {
		return
	}
	call.Args[0] = arg0

	conv.convert_method(call, "Sort")
}
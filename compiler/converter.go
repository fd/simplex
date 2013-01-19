package compiler

import (
	"github.com/fd/simplex/ast"
	"github.com/fd/simplex/token"
	"github.com/fd/simplex/types"
)

func (c *Context) convert_sx_to_go() error {

	for _, name := range c.SxFiles {
		file := c.AstFiles[name]

		runtime_name := ""
		for _, imp := range file.Imports {
			if imp.Path.Value == `"github.com/fd/simplex/runtime"` {
				if imp.Name == nil {
					runtime_name = "runtime"
				} else {
					runtime_name = imp.Name.Name
				}
			}
		}
		if runtime_name == "" {
			runtime_name = "sx_runtime"
			imp := &ast.ImportSpec{
				Name: ast.NewIdent("sx_runtime"),
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"github.com/fd/simplex/runtime"`,
				},
			}
			decl := &ast.GenDecl{
				Tok:   token.IMPORT,
				Specs: []ast.Spec{imp},
			}
			file.Decls = append([]ast.Decl{decl}, file.Decls...)
			file.Imports = append(file.Imports, imp)
		}

		ast.Walk(&builtin_function_conv{c.NodeTypes, runtime_name}, file)
	}

	return nil
}

type (
	builtin_function_conv struct {
		mapping      map[ast.Node]types.Type
		runtime_name string
	}
)

func (conv *builtin_function_conv) Visit(node ast.Node) ast.Visitor {
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
		if ident, ok := n.Fun.(*ast.Ident); ok && ident.Name == "make" {
			conv.convert_make(n)
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

	case *ast.FuncDecl:
		if n.Recv != nil {
			var (
				field    = n.Recv.List[0]
				typ_expr = field.Type
			)

			if ident, ok := typ_expr.(*ast.Ident); ok {
				var (
					typ  = conv.mapping[typ_expr]
					decl = ident.Obj.Decl
				)

				if named_typ, ok := typ.(*types.NamedType); ok {
					if _, ok := named_typ.Underlying.(types.Viewish); ok {
						if spec, ok := decl.(*ast.TypeSpec); ok {
							ast.Walk(conv, spec)

							var (
								interf = spec.Type.(*ast.InterfaceType)
								m      = interf.Methods.List
							)

							m = append(m, &ast.Field{
								Names: []*ast.Ident{n.Name},
								Type:  n.Type,
							})
							interf.Methods.List = m

							field.Type = ast.NewIdent("sx_" + type_name(named_typ.Underlying))
						}
					}
				}
			}
		}

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
			n.Type = &ast.InterfaceType{
				Methods: &ast.FieldList{
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
	return conv
}

func (conv *builtin_function_conv) convert_type(expr ast.Expr) ast.Expr {
	switch expr.(type) {
	case *ast.ViewType, *ast.TableType:
		typ := conv.mapping[expr]
		name := type_name(typ)
		return ast.NewIdent(name)
	}
	return expr
}

func (conv *builtin_function_conv) convert_make(call *ast.CallExpr) {
	typ := conv.mapping[call]
	if _, ok := typ.(*types.Table); !ok {
		return
	}

	call.Fun = ast.NewIdent("new_" + type_name(typ))
	call.Args = nil
}

func (conv *builtin_function_conv) convert_method(call *ast.CallExpr, name string) {
	var (
		recv  = call.Fun.(*ast.SelectorExpr).X
		args  = call.Args
		o_typ = conv.mapping[call]
	)

	args = append([]ast.Expr{recv}, args...)

	call.Fun = ast.NewIdent("wrap_" + type_name(o_typ))
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

func (conv *builtin_function_conv) convert_select(call *ast.CallExpr) {
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

	conv.convert_method(call, "Select")
}

func (conv *builtin_function_conv) convert_reject(call *ast.CallExpr) {
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

	conv.convert_method(call, "Reject")
}

func (conv *builtin_function_conv) convert_detect(call *ast.CallExpr) {
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

	conv.convert_method(call, "Detect")
}

func (conv *builtin_function_conv) convert_collect(call *ast.CallExpr) {
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

	conv.convert_method(call, "Group")
}

func (conv *builtin_function_conv) convert_index(call *ast.CallExpr) {
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

	conv.convert_method(call, "Index")
}

func (conv *builtin_function_conv) convert_sort(call *ast.CallExpr) {
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

	conv.convert_method(call, "Sort")
}

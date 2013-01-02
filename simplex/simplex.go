package main

import (
	"fmt"
	"github.com/fd/w/simplex/compiler"
)

func main() {
	_, err := compiler.ImportResolved(
		"github.com/fd/w/simplex/example", ".")
	if err != nil {
		fmt.Println(err)
	}
}

/*

undefined
nil
false
true
int
float
string
array
object

any

*/

/*
const smplx = `
package main

const Locations = Source.where(M._type == "location").sort(M.name)

func (s string) Lower () string {
  return undefined
}

func (i int) Add5 () int {
  source(int)
  return i + 5
}

func (i int.(view)) Add6 () int {
  return i + 5
}
`

func main() {
	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, "example.smplx", smplx, 0)
	if err != nil {
		fmt.Printf("error: %s", err)
	}


	ast.Walk(&replace_Expr{func(node ast.Expr) ast.Expr {
		if ident, ok := node.(*ast.Ident); ok && ident != nil && ident.Name == "undefined" {
			return &ast.CompositeLit{
				Type: &ast.Ident{Name: "Undefined"},
				Elts: []ast.Expr{&ast.BasicLit{
					Kind:  token.STRING,
					Value: strconv.Quote(fset.Position(node.Pos()).String()),
				}},
			}

		}
		return node
	}}, file)

	ast.Walk(&replace_Expr{func(node ast.Expr) ast.Expr {
		if binop, ok := node.(*ast.BinaryExpr); ok {
			op_name := ""
			switch binop.Op {
			case token.ADD: // +
        op_name = "ADD"
			case token.SUB: // -
        op_name = "SUB"
			case token.MUL: // *
        op_name = "MUL"
			case token.QUO: // /
        op_name = "QUO"
			case token.REM: // %
        op_name = "REM"

			case token.AND: // &
        op_name = "AND"
			case token.OR: // |
        op_name = "OR"
			case token.XOR: // ^
        op_name = "XOR"
			case token.SHL: // <<
        op_name = "SHL"
			case token.SHR: // >>
        op_name = "SHR"
			case token.AND_NOT: // &^
        op_name = "AND_NOT"

			case token.ADD_ASSIGN: // +=
        op_name = "ADD_ASSIGN"
			case token.SUB_ASSIGN: // -=
        op_name = "SUB_ASSIGN"
			case token.MUL_ASSIGN: // *=
        op_name = "MUL_ASSIGN"
			case token.QUO_ASSIGN: // /=
        op_name = "QUO_ASSIGN"
			case token.REM_ASSIGN: // %=
        op_name = "REM_ASSIGN"

			case token.AND_ASSIGN: // &=
        op_name = "AND_ASSIGN"
			case token.OR_ASSIGN: // |=
        op_name = "OR_ASSIGN"
			case token.XOR_ASSIGN: // ^=
        op_name = "XOR_ASSIGN"
			case token.SHL_ASSIGN: // <<=
        op_name = "SHL_ASSIGN"
			case token.SHR_ASSIGN: // >>=
        op_name = "SHR_ASSIGN"
			case token.AND_NOT_ASSIGN: // &^=
        op_name = "AND_NOT_ASSIGN"

			case token.LAND: // &&
        op_name = "LAND"
			case token.LOR: // ||
        op_name = "LOR"
			case token.ARROW: // <-
        op_name = "ARROW"
			case token.INC: // ++
        op_name = "INC"
			case token.DEC: // --
        op_name = "DEC"

			case token.EQL: // ==
        op_name = "EQL"
			case token.LSS: // <
        op_name = "LSS"
			case token.GTR: // >
        op_name = "GTR"
			case token.ASSIGN: // =
        op_name = "ASSIGN"
			case token.NOT: // !
        op_name = "NOT"

			case token.NEQ: // !=
        op_name = "NEQ"
			case token.LEQ: // <=
        op_name = "LEQ"
			case token.GEQ: // >=
        op_name = "GEQ"
			case token.DEFINE: // :=
        op_name = "DEFINE"
			case token.ELLIPSIS: // ...
        op_name = "ELLIPSIS"
			}
			return &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent("simplex_runtime"),
					Sel: ast.NewIdent("BINOP_" + op_name),
				},
				Args: []ast.Expr{binop.X, binop.Y},
			}
		}
		return node
	}}, file)

	ast.Walk(&replace_Type{func(node ast.Expr) ast.Expr {
		if ident, ok := node.(*ast.Ident); ok && ident != nil && ident.Name == "string" {
			return &ast.Ident{NamePos: ident.NamePos, Name: "String"}
		}
		return node
	}}, file)

	ast.Print(fset, file)
	c := printer.Config{printer.TabIndent | printer.SourcePos, 8}
	c.Fprint(os.Stdout, fset, file)
}

type replace_Type struct{ f func(node ast.Expr) ast.Expr }

func (v *replace_Type) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.ArrayType:
		n.Elt = v.Perform(n.Elt)
	case *ast.ChanType:
		n.Value = v.Perform(n.Value)
	case *ast.CompositeLit:
		n.Type = v.Perform(n.Type)
	case *ast.Ellipsis:
		n.Elt = v.Perform(n.Elt)
	case *ast.Field:
		n.Type = v.Perform(n.Type)
	case *ast.MapType:
		n.Key = v.Perform(n.Key)
		n.Value = v.Perform(n.Value)
	case *ast.ParenExpr:
		n.X = v.Perform(n.X)
	case *ast.SelectorExpr:
		n.X = v.Perform(n.X)
	case *ast.StarExpr:
		n.X = v.Perform(n.X)
	case *ast.SwitchStmt:
		n.Tag = v.Perform(n.Tag)
	case *ast.TypeAssertExpr:
		n.Type = v.Perform(n.Type)
	case *ast.TypeSpec:
		n.Type = v.Perform(n.Type)
	case *ast.ValueSpec:
		n.Type = v.Perform(n.Type)
	}
	return v
}
func (v *replace_Type) PerformN(nodes []ast.Expr) []ast.Expr {
	res := make([]ast.Expr, len(nodes))
	for i, node := range nodes {
		res[i] = v.Perform(node)
	}
	return res
}
func (v *replace_Type) Perform(node ast.Expr) ast.Expr {
	return v.f(node)
}

type replace_Expr struct{ f func(node ast.Expr) ast.Expr }

func (v *replace_Expr) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.ArrayType:
		n.Len = v.Perform(n.Len)
		n.Elt = v.Perform(n.Elt)
	case *ast.AssignStmt:
		n.Lhs = v.PerformN(n.Lhs)
		n.Rhs = v.PerformN(n.Rhs)
	case *ast.BinaryExpr:
		n.X = v.Perform(n.X)
		n.Y = v.Perform(n.Y)
	case *ast.CallExpr:
		n.Fun = v.Perform(n.Fun)
		n.Args = v.PerformN(n.Args)
	case *ast.CaseClause:
		n.List = v.PerformN(n.List)
	case *ast.ChanType:
		n.Value = v.Perform(n.Value)
	case *ast.CompositeLit:
		n.Type = v.Perform(n.Type)
		n.Elts = v.PerformN(n.Elts)
	case *ast.Ellipsis:
		n.Elt = v.Perform(n.Elt)
	case *ast.Field:
		n.Type = v.Perform(n.Type)
	case *ast.ForStmt:
		n.Cond = v.Perform(n.Cond)
	case *ast.IfStmt:
		n.Cond = v.Perform(n.Cond)
	case *ast.IncDecStmt:
		n.X = v.Perform(n.X)
	case *ast.IndexExpr:
		n.X = v.Perform(n.X)
		n.Index = v.Perform(n.Index)
	case *ast.KeyValueExpr:
		n.Key = v.Perform(n.Key)
		n.Value = v.Perform(n.Value)
	case *ast.MapType:
		n.Key = v.Perform(n.Key)
		n.Value = v.Perform(n.Value)
	case *ast.ParenExpr:
		n.X = v.Perform(n.X)
	case *ast.RangeStmt:
		n.Key = v.Perform(n.Key)
		n.Value = v.Perform(n.Value)
		n.X = v.Perform(n.X)
	case *ast.ReturnStmt:
		n.Results = v.PerformN(n.Results)
	case *ast.SelectorExpr:
		n.X = v.Perform(n.X)
	case *ast.SendStmt:
		n.Chan = v.Perform(n.Chan)
		n.Value = v.Perform(n.Value)
	case *ast.SliceExpr:
		n.X = v.Perform(n.X)
		n.High = v.Perform(n.High)
		n.Low = v.Perform(n.Low)
	case *ast.StarExpr:
		n.X = v.Perform(n.X)
	case *ast.SwitchStmt:
		n.Tag = v.Perform(n.Tag)
	case *ast.TypeAssertExpr:
		n.X = v.Perform(n.X)
		n.Type = v.Perform(n.Type)
	case *ast.TypeSpec:
		n.Type = v.Perform(n.Type)
	case *ast.UnaryExpr:
		n.X = v.Perform(n.X)
	case *ast.ValueSpec:
		n.Type = v.Perform(n.Type)
		n.Values = v.PerformN(n.Values)
	}
	return v
}
func (v *replace_Expr) PerformN(nodes []ast.Expr) []ast.Expr {
	res := make([]ast.Expr, len(nodes))
	for i, node := range nodes {
		res[i] = v.Perform(node)
	}
	return res
}
func (v *replace_Expr) Perform(node ast.Expr) ast.Expr {
	return v.f(node)
}

var (
	scope    *ast.Scope // current scope to use for initialization
	Universe *ast.Scope
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

func resolve(scope *ast.Scope, ident *ast.Ident) bool {
	for ; scope != nil; scope = scope.Outer {
		if obj := scope.Lookup(ident.Name); obj != nil {
			ident.Obj = obj
			return true
		}
	}
	return false
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

	defType("undefined")

	defConst("true")
	defConst("false")
	defConst("iota")
	defConst("nil")

	defConst("M")
	defConst("Source")

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

	defFun("where")
	defFun("sort")
}
*/

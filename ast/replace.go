package ast

import "fmt"

type Replacer interface {
	Replace(node Node) (w Replacer, n Node)
}

func replaceSelect(x, y Node) Node {
	if x != nil {
		return x
	}
	return y
}

func replaceIdentList(v Replacer, list []*Ident) {
	for i, x := range list {
		list[i] = Replace(v, x).(*Ident)
	}
}

func replaceExprList(v Replacer, list []Expr) {
	for i, x := range list {
		list[i] = Replace(v, x).(Expr)
	}
}

func replaceStmtList(v Replacer, list []Stmt) {
	for i, x := range list {
		list[i] = Replace(v, x).(Stmt)
	}
}

func replaceDeclList(v Replacer, list []Decl) {
	for i, x := range list {
		list[i] = Replace(v, x).(Decl)
	}
}

func Replace(v Replacer, node Node) Node {
	v, y := v.Replace(node)
	if y != nil {
		node = y
	}
	if v == nil {
		return node
	}

	// replace children
	// (the order of the cases matches the order
	// of the corresponding node types in ast.go)
	switch n := node.(type) {
	// Comments and fields
	case *Comment:
		// nothing to do

	case *CommentGroup:
		for i, c := range n.List {
			n.List[i] = Replace(v, c).(*Comment)
		}

	case *Field:
		if n.Doc != nil {
			n.Doc = Replace(v, n.Doc).(*CommentGroup)
		}
		replaceIdentList(v, n.Names)
		n.Type = Replace(v, n.Type).(Expr)
		if n.Tag != nil {
			n.Tag = Replace(v, n.Tag).(*BasicLit)
		}
		if n.Comment != nil {
			n.Comment = Replace(v, n.Comment).(*CommentGroup)
		}

	case *FieldList:
		for i, f := range n.List {
			n.List[i] = Replace(v, f).(*Field)
		}

	// Expressions
	case *BadExpr, *Ident, *BasicLit:
		// nothing to do

	case *Ellipsis:
		if n.Elt != nil {
			n.Elt = Replace(v, n.Elt).(Expr)
		}

	case *FuncLit:
		n.Type = Replace(v, n.Type).(*FuncType)
		n.Body = Replace(v, n.Body).(*BlockStmt)

	case *CompositeLit:
		if n.Type != nil {
			n.Type = Replace(v, n.Type).(Expr)
		}
		replaceExprList(v, n.Elts)

	case *ParenExpr:
		n.X = Replace(v, n.X).(Expr)

	case *SelectorExpr:
		n.X = Replace(v, n.X).(Expr)
		n.Sel = Replace(v, n.Sel).(*Ident)

	case *IndexExpr:
		n.X = Replace(v, n.X).(Expr)
		n.Index = Replace(v, n.Index).(Expr)

	case *SliceExpr:
		n.X = Replace(v, n.X).(Expr)
		if n.Low != nil {
			n.Low = Replace(v, n.Low).(Expr)
		}
		if n.High != nil {
			n.High = Replace(v, n.High).(Expr)
		}

	case *TypeAssertExpr:
		n.X = Replace(v, n.X).(Expr)
		if n.Type != nil {
			n.Type = Replace(v, n.Type).(Expr)
		}

	case *CallExpr:
		n.Fun = Replace(v, n.Fun).(Expr)
		replaceExprList(v, n.Args)

	case *StarExpr:
		n.X = Replace(v, n.X).(Expr)

	case *UnaryExpr:
		n.X = Replace(v, n.X).(Expr)

	case *BinaryExpr:
		n.X = Replace(v, n.X).(Expr)
		n.Y = Replace(v, n.Y).(Expr)

	case *KeyValueExpr:
		n.Key = Replace(v, n.Key).(Expr)
		n.Value = Replace(v, n.Value).(Expr)

	// Types
	case *ArrayType:
		if n.Len != nil {
			n.Len = Replace(v, n.Len).(Expr)
		}
		n.Elt = Replace(v, n.Elt).(Expr)

	case *StructType:
		n.Fields = Replace(v, n.Fields).(*FieldList)

	case *FuncType:
		if n.Params != nil {
			n.Params = Replace(v, n.Params).(*FieldList)
		}
		if n.Results != nil {
			n.Results = Replace(v, n.Results).(*FieldList)
		}

	case *InterfaceType:
		n.Methods = Replace(v, n.Methods).(*FieldList)

	case *MapType:
		n.Key = Replace(v, n.Key).(Expr)
		n.Value = Replace(v, n.Value).(Expr)

	case *ChanType:
		n.Value = Replace(v, n.Value).(Expr)

	// Statements
	case *BadStmt:
		// nothing to do

	case *DeclStmt:
		Replace(v, n.Decl)

	case *EmptyStmt:
		// nothing to do

	case *LabeledStmt:
		n.Label = Replace(v, n.Label).(*Ident)
		n.Stmt = Replace(v, n.Stmt).(Stmt)

	case *ExprStmt:
		n.X = Replace(v, n.X).(Expr)

	case *SendStmt:
		n.Chan = Replace(v, n.Chan).(Expr)
		n.Value = Replace(v, n.Value).(Expr)

	case *IncDecStmt:
		n.X = Replace(v, n.X).(Expr)

	case *AssignStmt:
		replaceExprList(v, n.Lhs)
		replaceExprList(v, n.Rhs)

	case *GoStmt:
		n.Call = Replace(v, n.Call).(*CallExpr)

	case *DeferStmt:
		n.Call = Replace(v, n.Call).(*CallExpr)

	case *ReturnStmt:
		replaceExprList(v, n.Results)

	case *BranchStmt:
		if n.Label != nil {
			n.Label = Replace(v, n.Label).(*Ident)
		}

	case *BlockStmt:
		replaceStmtList(v, n.List)

	case *IfStmt:
		if n.Init != nil {
			n.Init = Replace(v, n.Init).(Stmt)
		}
		n.Cond = Replace(v, n.Cond).(Expr)
		n.Body = Replace(v, n.Body).(*BlockStmt)
		if n.Else != nil {
			n.Else = Replace(v, n.Else).(Stmt)
		}

	case *CaseClause:
		replaceExprList(v, n.List)
		replaceStmtList(v, n.Body)

	case *SwitchStmt:
		if n.Init != nil {
			n.Init = Replace(v, n.Init).(Stmt)
		}
		if n.Tag != nil {
			n.Tag = Replace(v, n.Tag).(Expr)
		}
		n.Body = Replace(v, n.Body).(*BlockStmt)

	case *TypeSwitchStmt:
		if n.Init != nil {
			n.Init = Replace(v, n.Init).(Stmt)
		}
		n.Assign = Replace(v, n.Assign).(Stmt)
		n.Body = Replace(v, n.Body).(*BlockStmt)

	case *CommClause:
		if n.Comm != nil {
			n.Comm = Replace(v, n.Comm).(Stmt)
		}
		replaceStmtList(v, n.Body)

	case *SelectStmt:
		n.Body = Replace(v, n.Body).(*BlockStmt)

	case *ForStmt:
		if n.Init != nil {
			n.Init = Replace(v, n.Init).(Stmt)
		}
		if n.Cond != nil {
			n.Cond = Replace(v, n.Cond).(Expr)
		}
		if n.Post != nil {
			n.Post = Replace(v, n.Post).(Stmt)
		}
		n.Body = Replace(v, n.Body).(*BlockStmt)

	case *RangeStmt:
		n.Key = Replace(v, n.Key).(Expr)
		if n.Value != nil {
			n.Value = Replace(v, n.Value).(Expr)
		}
		n.X = Replace(v, n.X).(Expr)
		n.Body = Replace(v, n.Body).(*BlockStmt)

	// Declarations
	case *ImportSpec:
		if n.Doc != nil {
			n.Doc = Replace(v, n.Doc).(*CommentGroup)
		}
		if n.Name != nil {
			n.Name = Replace(v, n.Name).(*Ident)
		}
		n.Path = Replace(v, n.Path).(*BasicLit)
		if n.Comment != nil {
			n.Comment = Replace(v, n.Comment).(*CommentGroup)
		}

	case *ValueSpec:
		if n.Doc != nil {
			n.Doc = Replace(v, n.Doc).(*CommentGroup)
		}
		replaceIdentList(v, n.Names)
		if n.Type != nil {
			n.Type = Replace(v, n.Type).(Expr)
		}
		replaceExprList(v, n.Values)
		if n.Comment != nil {
			n.Comment = Replace(v, n.Comment).(*CommentGroup)
		}

	case *TypeSpec:
		if n.Doc != nil {
			n.Doc = Replace(v, n.Doc).(*CommentGroup)
		}
		n.Name = Replace(v, n.Name).(*Ident)
		n.Type = Replace(v, n.Type).(Expr)
		if n.Comment != nil {
			n.Comment = Replace(v, n.Comment).(*CommentGroup)
		}

	case *BadDecl:
		// nothing to do

	case *GenDecl:
		if n.Doc != nil {
			n.Doc = Replace(v, n.Doc).(*CommentGroup)
		}
		for i, s := range n.Specs {
			n.Specs[i] = Replace(v, s).(Spec)
		}

	case *FuncDecl:
		if n.Doc != nil {
			n.Doc = Replace(v, n.Doc).(*CommentGroup)
		}
		if n.Recv != nil {
			n.Recv = Replace(v, n.Recv).(*FieldList)
		}
		n.Name = Replace(v, n.Name).(*Ident)
		n.Type = Replace(v, n.Type).(*FuncType)
		if n.Body != nil {
			n.Body = Replace(v, n.Body).(*BlockStmt)
		}

	// Files and packages
	case *File:
		if n.Doc != nil {
			n.Doc = Replace(v, n.Doc).(*CommentGroup)
		}
		n.Name = Replace(v, n.Name).(*Ident)
		replaceDeclList(v, n.Decls)
		// don't replace n.Comments - they have been
		// visited already through the individual
		// nodes

	case *Package:
		for i, f := range n.Files {
			n.Files[i] = Replace(v, f).(*File)
		}

	case *ViewType:
		if n.Key != nil {
			n.Key = Replace(v, n.Key).(Expr)
		}
		n.Value = Replace(v, n.Value).(Expr)

	case *TableType:
		n.Key = Replace(v, n.Key).(Expr)
		n.Value = Replace(v, n.Value).(Expr)

	default:
		fmt.Printf("ast.Replace: unexpected node type %T", n)
		panic("ast.Replace")
	}

	return node
}

package ast

import (
	go_ast "go/ast"
	go_token "go/token"
)

// ----------------------------------------------------------------------------
// === Comment
func (c *Comment) GoComment() *go_ast.Comment {
	return &go_ast.Comment{c.Slash, c.Text}
}
func (c *Comment) GoNode() go_ast.Node { return c.GoComment() }

// ----------------------------------------------------------------------------
// === CommentGroup
func (g *CommentGroup) GoCommentGroup() *go_ast.CommentGroup {
	if g == nil {
		return nil
	}

	list := make([]*go_ast.Comment, len(g.List))
	for i, c := range g.List {
		list[i] = c.GoComment()
	}
	return &go_ast.CommentGroup{list}
}
func (g *CommentGroup) GoNode() go_ast.Node { return g.GoCommentGroup() }

// ----------------------------------------------------------------------------
// === Field
func (f *Field) GoField() *go_ast.Field {
	if f == nil {
		return nil
	}

	var names []*go_ast.Ident
	if f.Names != nil {
		names = make([]*go_ast.Ident, len(f.Names))
		for i, ident := range f.Names {
			names[i] = ident.GoIdent()
		}
	}
	return &go_ast.Field{
		f.Doc.GoCommentGroup(),
		names,
		f.Type.GoExpr(),
		f.Tag.GoBasicLit(),
		f.Comment.GoCommentGroup(),
	}
}
func (f *Field) GoNode() go_ast.Node { return f.GoField() }

// ----------------------------------------------------------------------------
// === FieldList
func (f *FieldList) GoFieldList() *go_ast.FieldList {
	if f == nil {
		return nil
	}

	var list []*go_ast.Field
	if f.List != nil {
		list = make([]*go_ast.Field, len(f.List))
		for i, f := range f.List {
			list[i] = f.GoField()
		}
	}

	return &go_ast.FieldList{
		f.Opening,
		list,
		f.Closing,
	}
}
func (f *FieldList) GoNode() go_ast.Node { return f.GoFieldList() }

// ----------------------------------------------------------------------------
// === BadExpr
func (e *BadExpr) GoBadExpr() *go_ast.BadExpr {
	if e == nil {
		return nil
	}
	return &go_ast.BadExpr{e.From, e.To}
}
func (e *BadExpr) GoExpr() go_ast.Expr { return e.GoBadExpr() }
func (e *BadExpr) GoNode() go_ast.Node { return e.GoBadExpr() }

// ----------------------------------------------------------------------------
// === Ident
func (e *Ident) GoIdent() *go_ast.Ident {
	if e == nil {
		return nil
	}
	return &go_ast.Ident{e.NamePos, e.Name, e.Obj}
}
func (e *Ident) GoExpr() go_ast.Expr { return e.GoIdent() }
func (e *Ident) GoNode() go_ast.Node { return e.GoIdent() }

// ----------------------------------------------------------------------------
// === Ellipsis
func (e *Ellipsis) GoEllipsis() *go_ast.Ellipsis {
	if e == nil {
		return nil
	}
	return &go_ast.Ellipsis{e.Ellipsis, e.Elt.GoExpr()}
}
func (e *Ellipsis) GoExpr() go_ast.Expr { return e.GoEllipsis() }
func (e *Ellipsis) GoNode() go_ast.Node { return e.GoEllipsis() }

// ----------------------------------------------------------------------------
// === BasicLit
func (e *BasicLit) GoBasicLit() *go_ast.BasicLit {
	if e == nil {
		return nil
	}
	return &go_ast.BasicLit{e.ValuePos, go_token.Token(e.Kind), e.Value}
}
func (e *BasicLit) GoExpr() go_ast.Expr { return e.GoBasicLit() }
func (e *BasicLit) GoNode() go_ast.Node { return e.GoBasicLit() }

// ----------------------------------------------------------------------------
// === FuncLit
func (e *FuncLit) GoFuncLit() *go_ast.FuncLit {
	if e == nil {
		return nil
	}
	return &go_ast.FuncLit{e.Type.GoFuncType(), e.Body.GoBlockStmt()}
}
func (e *FuncLit) GoExpr() go_ast.Expr { return e.GoFuncLit() }
func (e *FuncLit) GoNode() go_ast.Node { return e.GoFuncLit() }

// ----------------------------------------------------------------------------
// === CompositeLit
func (e *CompositeLit) GoCompositeLit() *go_ast.CompositeLit {
	if e == nil {
		return nil
	}

	var elts []go_ast.Expr
	if e.Elts != nil {
		elts = make([]go_ast.Expr, len(e.Elts))
		for i, f := range e.Elts {
			elts[i] = f.GoExpr()
		}
	}

	return &go_ast.CompositeLit{e.Type.GoExpr(), e.Lbrace, elts, e.Rbrace}
}
func (e *CompositeLit) GoExpr() go_ast.Expr { return e.GoCompositeLit() }
func (e *CompositeLit) GoNode() go_ast.Node { return e.GoCompositeLit() }

// ----------------------------------------------------------------------------
// === ParenExpr
func (e *ParenExpr) GoParenExpr() *go_ast.ParenExpr {
	if e == nil {
		return nil
	}

	return &go_ast.ParenExpr{e.Lparen, e.X.GoExpr(), e.Rparen}
}
func (e *ParenExpr) GoExpr() go_ast.Expr { return e.GoParenExpr() }
func (e *ParenExpr) GoNode() go_ast.Node { return e.GoParenExpr() }

// ----------------------------------------------------------------------------
// === SelectorExpr
func (e *SelectorExpr) GoSelectorExpr() *go_ast.SelectorExpr {
	if e == nil {
		return nil
	}

	return &go_ast.SelectorExpr{e.X.GoExpr(), e.Sel.GoIdent()}
}
func (e *SelectorExpr) GoExpr() go_ast.Expr { return e.GoSelectorExpr() }
func (e *SelectorExpr) GoNode() go_ast.Node { return e.GoSelectorExpr() }

// ----------------------------------------------------------------------------
// === IndexExpr
func (e *IndexExpr) GoIndexExpr() *go_ast.IndexExpr {
	if e == nil {
		return nil
	}

	return &go_ast.IndexExpr{e.X.GoExpr(), e.Lbrack, e.Index.GoExpr(), e.Rbrack}
}
func (e *IndexExpr) GoExpr() go_ast.Expr { return e.GoIndexExpr() }
func (e *IndexExpr) GoNode() go_ast.Node { return e.GoIndexExpr() }

// ----------------------------------------------------------------------------
// === SliceExpr
func (e *SliceExpr) GoSliceExpr() *go_ast.SliceExpr {
	if e == nil {
		return nil
	}

	return &go_ast.SliceExpr{
		e.X.GoExpr(),
		e.Lbrack,
		e.Low.GoExpr(),
		e.High.GoExpr(),
		e.Rbrack,
	}
}
func (e *SliceExpr) GoExpr() go_ast.Expr { return e.GoSliceExpr() }
func (e *SliceExpr) GoNode() go_ast.Node { return e.GoSliceExpr() }

// ----------------------------------------------------------------------------
// === TypeAssertExpr
func (e *TypeAssertExpr) GoTypeAssertExpr() *go_ast.TypeAssertExpr {
	if e == nil {
		return nil
	}

	return &go_ast.TypeAssertExpr{
		e.X.GoExpr(),
		e.Type.GoExpr(),
	}
}
func (e *TypeAssertExpr) GoExpr() go_ast.Expr { return e.GoTypeAssertExpr() }
func (e *TypeAssertExpr) GoNode() go_ast.Node { return e.GoTypeAssertExpr() }

// ----------------------------------------------------------------------------
// === CallExpr
func (e *CallExpr) GoCallExpr() *go_ast.CallExpr {
	if e == nil {
		return nil
	}

	var args []go_ast.Expr
	if e.Args != nil {
		args = make([]go_ast.Expr, len(e.Args))
		for i, f := range e.Args {
			args[i] = f.GoExpr()
		}
	}

	return &go_ast.CallExpr{
		e.Fun.GoExpr(),
		e.Lparen,
		args,
		e.Ellipsis,
		e.Rparen,
	}
}
func (e *CallExpr) GoExpr() go_ast.Expr { return e.GoCallExpr() }
func (e *CallExpr) GoNode() go_ast.Node { return e.GoCallExpr() }

// ----------------------------------------------------------------------------
// === StarExpr
func (e *StarExpr) GoStarExpr() *go_ast.StarExpr {
	if e == nil {
		return nil
	}

	return &go_ast.StarExpr{
		e.Star,
		e.X.GoExpr(),
	}
}
func (e *StarExpr) GoExpr() go_ast.Expr { return e.GoStarExpr() }
func (e *StarExpr) GoNode() go_ast.Node { return e.GoStarExpr() }

// ----------------------------------------------------------------------------
// === UnaryExpr
func (e *UnaryExpr) GoUnaryExpr() *go_ast.UnaryExpr {
	if e == nil {
		return nil
	}

	return &go_ast.UnaryExpr{
		e.OpPos,
		go_token.Token(e.Op),
		e.X.GoExpr(),
	}
}
func (e *UnaryExpr) GoExpr() go_ast.Expr { return e.GoUnaryExpr() }
func (e *UnaryExpr) GoNode() go_ast.Node { return e.GoUnaryExpr() }

// ----------------------------------------------------------------------------
// === BinaryExpr
func (e *BinaryExpr) GoBinaryExpr() *go_ast.BinaryExpr {
	if e == nil {
		return nil
	}

	return &go_ast.BinaryExpr{
		e.X.GoExpr(),
		e.OpPos,
		go_token.Token(e.Op),
		e.Y.GoExpr(),
	}
}
func (e *BinaryExpr) GoExpr() go_ast.Expr { return e.GoBinaryExpr() }
func (e *BinaryExpr) GoNode() go_ast.Node { return e.GoBinaryExpr() }

// ----------------------------------------------------------------------------
// === KeyValueExpr
func (e *KeyValueExpr) GoKeyValueExpr() *go_ast.KeyValueExpr {
	if e == nil {
		return nil
	}

	return &go_ast.KeyValueExpr{
		e.Key.GoExpr(),
		e.Colon,
		e.Value.GoExpr(),
	}
}
func (e *KeyValueExpr) GoExpr() go_ast.Expr { return e.GoKeyValueExpr() }
func (e *KeyValueExpr) GoNode() go_ast.Node { return e.GoKeyValueExpr() }

// ----------------------------------------------------------------------------
// === ArrayType
func (e *ArrayType) GoArrayType() *go_ast.ArrayType {
	if e == nil {
		return nil
	}

	return &go_ast.ArrayType{
		e.Lbrack,
		e.Len.GoExpr(),
		e.Elt.GoExpr(),
	}
}
func (e *ArrayType) GoExpr() go_ast.Expr { return e.GoArrayType() }
func (e *ArrayType) GoNode() go_ast.Node { return e.GoArrayType() }

// ----------------------------------------------------------------------------
// === StructType
func (e *StructType) GoStructType() *go_ast.StructType {
	if e == nil {
		return nil
	}

	return &go_ast.StructType{
		e.Struct,
		e.Fields.GoFieldList(),
		e.Incomplete,
	}
}
func (e *StructType) GoExpr() go_ast.Expr { return e.GoStructType() }
func (e *StructType) GoNode() go_ast.Node { return e.GoStructType() }

// ----------------------------------------------------------------------------
// === FuncType
func (e *FuncType) GoFuncType() *go_ast.FuncType {
	if e == nil {
		return nil
	}

	return &go_ast.FuncType{
		e.Func,
		e.Params.GoFieldList(),
		e.Results.GoFieldList(),
	}
}
func (e *FuncType) GoExpr() go_ast.Expr { return e.GoFuncType() }
func (e *FuncType) GoNode() go_ast.Node { return e.GoFuncType() }

// ----------------------------------------------------------------------------
// === InterfaceType
func (e *InterfaceType) GoInterfaceType() *go_ast.InterfaceType {
	if e == nil {
		return nil
	}

	return &go_ast.InterfaceType{
		e.Interface,
		e.Methods.GoFieldList(),
		e.Incomplete,
	}
}
func (e *InterfaceType) GoExpr() go_ast.Expr { return e.GoInterfaceType() }
func (e *InterfaceType) GoNode() go_ast.Node { return e.GoInterfaceType() }

// ----------------------------------------------------------------------------
// === MapType
func (e *MapType) GoMapType() *go_ast.MapType {
	if e == nil {
		return nil
	}

	return &go_ast.MapType{
		e.Map,
		e.Key.GoExpr(),
		e.Value.GoExpr(),
	}
}
func (e *MapType) GoExpr() go_ast.Expr { return e.GoMapType() }
func (e *MapType) GoNode() go_ast.Node { return e.GoMapType() }

// ----------------------------------------------------------------------------
// === ViewType
func (e *ViewType) GoMapType() go_ast.Expr {
	if e == nil {
		return nil
	}

	panic("implement ViewType")
}
func (e *ViewType) GoExpr() go_ast.Expr { return e.GoMapType() }
func (e *ViewType) GoNode() go_ast.Node { return e.GoMapType() }

// ----------------------------------------------------------------------------
// === ChanType
func (e *ChanType) GoChanType() *go_ast.ChanType {
	if e == nil {
		return nil
	}

	return &go_ast.ChanType{
		e.Begin,
		go_ast.ChanDir(e.Dir),
		e.Value.GoExpr(),
	}
}
func (e *ChanType) GoExpr() go_ast.Expr { return e.GoChanType() }
func (e *ChanType) GoNode() go_ast.Node { return e.GoChanType() }

// ----------------------------------------------------------------------------
// === BadStmt
func (e *BadStmt) GoBadStmt() *go_ast.BadStmt {
	if e == nil {
		return nil
	}

	return &go_ast.BadStmt{
		e.From,
		e.To,
	}
}
func (e *BadStmt) GoStmt() go_ast.Stmt { return e.GoBadStmt() }
func (e *BadStmt) GoNode() go_ast.Node { return e.GoBadStmt() }

// ----------------------------------------------------------------------------
// === DeclStmt
func (e *DeclStmt) GoDeclStmt() *go_ast.DeclStmt {
	if e == nil {
		return nil
	}

	return &go_ast.DeclStmt{
		e.Decl.GoDecl(),
	}
}
func (e *DeclStmt) GoStmt() go_ast.Stmt { return e.GoDeclStmt() }
func (e *DeclStmt) GoNode() go_ast.Node { return e.GoDeclStmt() }

// ----------------------------------------------------------------------------
// === EmptyStmt
func (e *EmptyStmt) GoEmptyStmt() *go_ast.EmptyStmt {
	if e == nil {
		return nil
	}

	return &go_ast.EmptyStmt{
		e.Semicolon,
	}
}
func (e *EmptyStmt) GoStmt() go_ast.Stmt { return e.GoEmptyStmt() }
func (e *EmptyStmt) GoNode() go_ast.Node { return e.GoEmptyStmt() }

// ----------------------------------------------------------------------------
// === LabeledStmt
func (e *LabeledStmt) GoLabeledStmt() *go_ast.LabeledStmt {
	if e == nil {
		return nil
	}

	return &go_ast.LabeledStmt{
		e.Label.GoIdent(),
		e.Colon,
		e.Stmt.GoStmt(),
	}
}
func (e *LabeledStmt) GoStmt() go_ast.Stmt { return e.GoLabeledStmt() }
func (e *LabeledStmt) GoNode() go_ast.Node { return e.GoLabeledStmt() }

// ----------------------------------------------------------------------------
// === ExprStmt
func (e *ExprStmt) GoExprStmt() *go_ast.ExprStmt {
	if e == nil {
		return nil
	}

	return &go_ast.ExprStmt{
		e.X.GoExpr(),
	}
}
func (e *ExprStmt) GoStmt() go_ast.Stmt { return e.GoExprStmt() }
func (e *ExprStmt) GoNode() go_ast.Node { return e.GoExprStmt() }

// ----------------------------------------------------------------------------
// === SendStmt
func (e *SendStmt) GoSendStmt() *go_ast.SendStmt {
	if e == nil {
		return nil
	}

	return &go_ast.SendStmt{
		e.Chan.GoExpr(),
		e.Arrow,
		e.Value.GoExpr(),
	}
}
func (e *SendStmt) GoStmt() go_ast.Stmt { return e.GoSendStmt() }
func (e *SendStmt) GoNode() go_ast.Node { return e.GoSendStmt() }

// ----------------------------------------------------------------------------
// === IncDecStmt
func (e *IncDecStmt) GoIncDecStmt() *go_ast.IncDecStmt {
	if e == nil {
		return nil
	}

	return &go_ast.IncDecStmt{
		e.X.GoExpr(),
		e.TokPos,
		go_token.Token(e.Tok),
	}
}
func (e *IncDecStmt) GoStmt() go_ast.Stmt { return e.GoIncDecStmt() }
func (e *IncDecStmt) GoNode() go_ast.Node { return e.GoIncDecStmt() }

// ----------------------------------------------------------------------------
// === AssignStmt
func (e *AssignStmt) GoAssignStmt() *go_ast.AssignStmt {
	if e == nil {
		return nil
	}

	var lhs []go_ast.Expr
	if e.Lhs != nil {
		lhs = make([]go_ast.Expr, len(e.Lhs))
		for i, f := range e.Lhs {
			lhs[i] = f.GoExpr()
		}
	}

	var rhs []go_ast.Expr
	if e.Rhs != nil {
		rhs = make([]go_ast.Expr, len(e.Rhs))
		for i, f := range e.Rhs {
			rhs[i] = f.GoExpr()
		}
	}

	return &go_ast.AssignStmt{
		lhs,
		e.TokPos,
		go_token.Token(e.Tok),
		rhs,
	}
}
func (e *AssignStmt) GoStmt() go_ast.Stmt { return e.GoAssignStmt() }
func (e *AssignStmt) GoNode() go_ast.Node { return e.GoAssignStmt() }

// ----------------------------------------------------------------------------
// === GoStmt
func (e *GoStmt) GoGoStmt() *go_ast.GoStmt {
	if e == nil {
		return nil
	}

	return &go_ast.GoStmt{
		e.Go,
		e.Call.GoCallExpr(),
	}
}
func (e *GoStmt) GoStmt() go_ast.Stmt { return e.GoGoStmt() }
func (e *GoStmt) GoNode() go_ast.Node { return e.GoGoStmt() }

// ----------------------------------------------------------------------------
// === DeferStmt
func (e *DeferStmt) GoDeferStmt() *go_ast.DeferStmt {
	if e == nil {
		return nil
	}

	return &go_ast.DeferStmt{
		e.Defer,
		e.Call.GoCallExpr(),
	}
}
func (e *DeferStmt) GoStmt() go_ast.Stmt { return e.GoDeferStmt() }
func (e *DeferStmt) GoNode() go_ast.Node { return e.GoDeferStmt() }

// ----------------------------------------------------------------------------
// === ReturnStmt
func (e *ReturnStmt) GoReturnStmt() *go_ast.ReturnStmt {
	if e == nil {
		return nil
	}

	var results []go_ast.Expr
	if e.Results != nil {
		results = make([]go_ast.Expr, len(e.Results))
		for i, f := range e.Results {
			results[i] = f.GoExpr()
		}
	}

	return &go_ast.ReturnStmt{
		e.Return,
		results,
	}
}
func (e *ReturnStmt) GoStmt() go_ast.Stmt { return e.GoReturnStmt() }
func (e *ReturnStmt) GoNode() go_ast.Node { return e.GoReturnStmt() }

// ----------------------------------------------------------------------------
// === BranchStmt
func (e *BranchStmt) GoBranchStmt() *go_ast.BranchStmt {
	if e == nil {
		return nil
	}

	return &go_ast.BranchStmt{
		e.TokPos,
		go_token.Token(e.Tok),
		e.Label.GoIdent(),
	}
}
func (e *BranchStmt) GoStmt() go_ast.Stmt { return e.GoBranchStmt() }
func (e *BranchStmt) GoNode() go_ast.Node { return e.GoBranchStmt() }

// ----------------------------------------------------------------------------
// === BlockStmt
func (e *BlockStmt) GoBlockStmt() *go_ast.BlockStmt {
	if e == nil {
		return nil
	}

	var list []go_ast.Stmt
	if e.List != nil {
		list = make([]go_ast.Stmt, len(e.List))
		for i, f := range e.List {
			list[i] = f.GoStmt()
		}
	}

	return &go_ast.BlockStmt{
		e.Lbrace,
		list,
		e.Rbrace,
	}
}
func (e *BlockStmt) GoStmt() go_ast.Stmt { return e.GoBlockStmt() }
func (e *BlockStmt) GoNode() go_ast.Node { return e.GoBlockStmt() }

// ----------------------------------------------------------------------------
// === IfStmt
func (e *IfStmt) GoIfStmt() *go_ast.IfStmt {
	if e == nil {
		return nil
	}

	return &go_ast.IfStmt{
		e.If,
		e.Init.GoStmt(),
		e.Cond.GoExpr(),
		e.Body.GoBlockStmt(),
		e.Else.GoStmt(),
	}
}
func (e *IfStmt) GoStmt() go_ast.Stmt { return e.GoIfStmt() }
func (e *IfStmt) GoNode() go_ast.Node { return e.GoIfStmt() }

// ----------------------------------------------------------------------------
// === CaseClause
func (e *CaseClause) GoCaseClause() *go_ast.CaseClause {
	if e == nil {
		return nil
	}

	var list []go_ast.Expr
	if e.List != nil {
		list = make([]go_ast.Expr, len(e.List))
		for i, f := range e.List {
			list[i] = f.GoExpr()
		}
	}

	var body []go_ast.Stmt
	if e.Body != nil {
		body = make([]go_ast.Stmt, len(e.Body))
		for i, f := range e.Body {
			body[i] = f.GoStmt()
		}
	}

	return &go_ast.CaseClause{
		e.Case,
		list,
		e.Colon,
		body,
	}
}
func (e *CaseClause) GoStmt() go_ast.Stmt { return e.GoCaseClause() }
func (e *CaseClause) GoNode() go_ast.Node { return e.GoCaseClause() }

// ----------------------------------------------------------------------------
// === SwitchStmt
func (e *SwitchStmt) GoSwitchStmt() *go_ast.SwitchStmt {
	if e == nil {
		return nil
	}

	return &go_ast.SwitchStmt{
		e.Switch,
		e.Init.GoStmt(),
		e.Tag.GoExpr(),
		e.Body.GoBlockStmt(),
	}
}
func (e *SwitchStmt) GoStmt() go_ast.Stmt { return e.GoSwitchStmt() }
func (e *SwitchStmt) GoNode() go_ast.Node { return e.GoSwitchStmt() }

// ----------------------------------------------------------------------------
// === TypeSwitchStmt
func (e *TypeSwitchStmt) GoTypeSwitchStmt() *go_ast.TypeSwitchStmt {
	if e == nil {
		return nil
	}

	return &go_ast.TypeSwitchStmt{
		e.Switch,
		e.Init.GoStmt(),
		e.Assign.GoStmt(),
		e.Body.GoBlockStmt(),
	}
}
func (e *TypeSwitchStmt) GoStmt() go_ast.Stmt { return e.GoTypeSwitchStmt() }
func (e *TypeSwitchStmt) GoNode() go_ast.Node { return e.GoTypeSwitchStmt() }

// ----------------------------------------------------------------------------
// === CommClause
func (e *CommClause) GoCommClause() *go_ast.CommClause {
	if e == nil {
		return nil
	}

	var body []go_ast.Stmt
	if e.Body != nil {
		body = make([]go_ast.Stmt, len(e.Body))
		for i, f := range e.Body {
			body[i] = f.GoStmt()
		}
	}

	return &go_ast.CommClause{
		e.Case,
		e.Comm.GoStmt(),
		e.Colon,
		body,
	}
}
func (e *CommClause) GoStmt() go_ast.Stmt { return e.GoCommClause() }
func (e *CommClause) GoNode() go_ast.Node { return e.GoCommClause() }

// ----------------------------------------------------------------------------
// === SelectStmt
func (e *SelectStmt) GoSelectStmt() *go_ast.SelectStmt {
	if e == nil {
		return nil
	}

	return &go_ast.SelectStmt{
		e.Select,
		e.Body.GoBlockStmt(),
	}
}
func (e *SelectStmt) GoStmt() go_ast.Stmt { return e.GoSelectStmt() }
func (e *SelectStmt) GoNode() go_ast.Node { return e.GoSelectStmt() }

// ----------------------------------------------------------------------------
// === ForStmt
func (e *ForStmt) GoForStmt() *go_ast.ForStmt {
	if e == nil {
		return nil
	}

	return &go_ast.ForStmt{
		e.For,
		e.Init.GoStmt(),
		e.Cond.GoExpr(),
		e.Post.GoStmt(),
		e.Body.GoBlockStmt(),
	}
}
func (e *ForStmt) GoStmt() go_ast.Stmt { return e.GoForStmt() }
func (e *ForStmt) GoNode() go_ast.Node { return e.GoForStmt() }

// ----------------------------------------------------------------------------
// === RangeStmt
func (e *RangeStmt) GoRangeStmt() *go_ast.RangeStmt {
	if e == nil {
		return nil
	}

	return &go_ast.RangeStmt{
		e.For,
		e.Key.GoExpr(),
		e.Value.GoExpr(),
		e.TokPos,
		go_token.Token(e.Tok),
		e.X.GoExpr(),
		e.Body.GoBlockStmt(),
	}
}
func (e *RangeStmt) GoStmt() go_ast.Stmt { return e.GoRangeStmt() }
func (e *RangeStmt) GoNode() go_ast.Node { return e.GoRangeStmt() }

// ----------------------------------------------------------------------------
// === ImportSpec
func (e *ImportSpec) GoImportSpec() *go_ast.ImportSpec {
	if e == nil {
		return nil
	}

	return &go_ast.ImportSpec{
		e.Doc.GoCommentGroup(),
		e.Name.GoIdent(),
		e.Path.GoBasicLit(),
		e.Comment.GoCommentGroup(),
		e.EndPos,
	}
}
func (e *ImportSpec) GoSpec() go_ast.Spec { return e.GoImportSpec() }
func (e *ImportSpec) GoNode() go_ast.Node { return e.GoImportSpec() }

// ----------------------------------------------------------------------------
// === ValueSpec
func (e *ValueSpec) GoValueSpec() *go_ast.ValueSpec {
	if e == nil {
		return nil
	}

	var names []*go_ast.Ident
	if e.Names != nil {
		names = make([]*go_ast.Ident, len(e.Names))
		for i, f := range e.Names {
			names[i] = f.GoIdent()
		}
	}

	var values []go_ast.Expr
	if e.Values != nil {
		values = make([]go_ast.Expr, len(e.Values))
		for i, f := range e.Values {
			values[i] = f.GoExpr()
		}
	}

	return &go_ast.ValueSpec{
		e.Doc.GoCommentGroup(),
		names,
		e.Type.GoExpr(),
		values,
		e.Comment.GoCommentGroup(),
	}
}
func (e *ValueSpec) GoSpec() go_ast.Spec { return e.GoValueSpec() }
func (e *ValueSpec) GoNode() go_ast.Node { return e.GoValueSpec() }

// ----------------------------------------------------------------------------
// === TypeSpec
func (e *TypeSpec) GoTypeSpec() *go_ast.TypeSpec {
	if e == nil {
		return nil
	}

	return &go_ast.TypeSpec{
		e.Doc.GoCommentGroup(),
		e.Name.GoIdent(),
		e.Type.GoExpr(),
		e.Comment.GoCommentGroup(),
	}
}
func (e *TypeSpec) GoSpec() go_ast.Spec { return e.GoTypeSpec() }
func (e *TypeSpec) GoNode() go_ast.Node { return e.GoTypeSpec() }

// ----------------------------------------------------------------------------
// === BadDecl
func (e *BadDecl) GoBadDecl() *go_ast.BadDecl {
	if e == nil {
		return nil
	}

	return &go_ast.BadDecl{
		e.From,
		e.To,
	}
}
func (e *BadDecl) GoDecl() go_ast.Decl { return e.GoBadDecl() }
func (e *BadDecl) GoNode() go_ast.Node { return e.GoBadDecl() }

// ----------------------------------------------------------------------------
// === GenDecl
func (e *GenDecl) GoGenDecl() *go_ast.GenDecl {
	if e == nil {
		return nil
	}

	var specs []go_ast.Spec
	if e.Specs != nil {
		specs = make([]go_ast.Spec, len(e.Specs))
		for i, f := range e.Specs {
			specs[i] = f.GoSpec()
		}
	}

	return &go_ast.GenDecl{
		e.Doc.GoCommentGroup(),
		e.TokPos,
		go_token.Token(e.Tok),
		e.Lparen,
		specs,
		e.Rparen,
	}
}
func (e *GenDecl) GoDecl() go_ast.Decl { return e.GoGenDecl() }
func (e *GenDecl) GoNode() go_ast.Node { return e.GoGenDecl() }

// ----------------------------------------------------------------------------
// === FuncDecl
func (e *FuncDecl) GoFuncDecl() *go_ast.FuncDecl {
	if e == nil {
		return nil
	}

	return &go_ast.FuncDecl{
		e.Doc.GoCommentGroup(),
		e.Recv.GoFieldList(),
		e.Name.GoIdent(),
		e.Type.GoFuncType(),
		e.Body.GoBlockStmt(),
	}
}
func (e *FuncDecl) GoDecl() go_ast.Decl { return e.GoFuncDecl() }
func (e *FuncDecl) GoNode() go_ast.Node { return e.GoFuncDecl() }

// ----------------------------------------------------------------------------
// === File
func (e *File) GoFile() *go_ast.File {
	if e == nil {
		return nil
	}

	var decls []go_ast.Decl
	if e.Decls != nil {
		decls = make([]go_ast.Decl, len(e.Decls))
		for i, f := range e.Decls {
			decls[i] = f.GoDecl()
		}
	}

	var imports []*go_ast.ImportSpec
	if e.Imports != nil {
		imports = make([]*go_ast.ImportSpec, len(e.Imports))
		for i, f := range e.Imports {
			imports[i] = f.GoImportSpec()
		}
	}

	var unresolved []*go_ast.Ident
	if e.Unresolved != nil {
		unresolved = make([]*go_ast.Ident, len(e.Unresolved))
		for i, f := range e.Unresolved {
			unresolved[i] = f.GoIdent()
		}
	}

	var comments []*go_ast.CommentGroup
	if e.Comments != nil {
		comments = make([]*go_ast.CommentGroup, len(e.Comments))
		for i, f := range e.Comments {
			comments[i] = f.GoCommentGroup()
		}
	}

	return &go_ast.File{
		e.Doc.GoCommentGroup(),
		e.Package,
		e.Name.GoIdent(),
		decls,
		e.Scope,
		imports,
		unresolved,
		comments,
	}
}
func (e *File) GoNode() go_ast.Node { return e.GoFile() }

// ----------------------------------------------------------------------------
// === Package
func (e *Package) GoPackage() *go_ast.Package {
	if e == nil {
		return nil
	}

	var files map[string]*go_ast.File
	if e.Files != nil {
		files = make(map[string]*go_ast.File, len(e.Files))
		for n, f := range e.Files {
			files[n] = f.GoFile()
		}
	}

	return &go_ast.Package{
		e.Name,
		e.Scope,
		e.Imports,
		files,
	}
}
func (e *Package) GoNode() go_ast.Node { return e.GoPackage() }

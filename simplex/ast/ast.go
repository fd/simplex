// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ast declares the types used to represent syntax trees for Go
// packages.
//
package ast

import (
	"github.com/fd/w/simplex/token"
	go_ast "go/ast"
	go_token "go/token"
	"strings"
	"unicode"
	"unicode/utf8"
)

// ----------------------------------------------------------------------------
// Interfaces
//
// There are 3 main classes of nodes: Expressions and type nodes,
// statement nodes, and declaration nodes. The node names usually
// match the corresponding Go spec production names to which they
// correspond. The node fields correspond to the individual parts
// of the respective productions.
//
// All nodes contain position information marking the beginning of
// the corresponding source text segment; it is accessible via the
// Pos accessor method. Nodes may contain additional position info
// for language constructs where comments may be found between parts
// of the construct (typically any larger, parenthesized subpart).
// That position information is needed to properly position comments
// when printing the construct.

// All node types implement the Node interface.
type Node interface {
	Pos() go_token.Pos // position of first character belonging to the node
	End() go_token.Pos // position of first character immediately after the node
	GoNode() go_ast.Node
}

// All expression nodes implement the Expr interface.
type Expr interface {
	Node
	GoExpr() go_ast.Expr
	exprNode()
}

// All statement nodes implement the Stmt interface.
type Stmt interface {
	Node
	GoStmt() go_ast.Stmt
	stmtNode()
}

// All declaration nodes implement the Decl interface.
type Decl interface {
	Node
	GoDecl() go_ast.Decl
	declNode()
}

// ----------------------------------------------------------------------------
// Comments

// A Comment node represents a single //-style or /*-style comment.
type Comment struct {
	Slash go_token.Pos // position of "/" starting the comment
	Text  string       // comment text (excluding '\n' for //-style comments)
}

func (c *Comment) Pos() go_token.Pos { return c.Slash }
func (c *Comment) End() go_token.Pos { return go_token.Pos(int(c.Slash) + len(c.Text)) }

// A CommentGroup represents a sequence of comments
// with no other tokens and no empty lines between.
//
type CommentGroup struct {
	List []*Comment // len(List) > 0
}

func (g *CommentGroup) Pos() go_token.Pos { return g.List[0].Pos() }
func (g *CommentGroup) End() go_token.Pos { return g.List[len(g.List)-1].End() }

func isWhitespace(ch byte) bool { return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' }

func stripTrailingWhitespace(s string) string {
	i := len(s)
	for i > 0 && isWhitespace(s[i-1]) {
		i--
	}
	return s[0:i]
}

// Text returns the text of the comment.
// Comment markers (//, /*, and */), the first space of a line comment, and
// leading and trailing empty lines are removed. Multiple empty lines are
// reduced to one, and trailing space on lines is trimmed. Unless the result
// is empty, it is newline-terminated.
//
func (g *CommentGroup) Text() string {
	if g == nil {
		return ""
	}
	comments := make([]string, len(g.List))
	for i, c := range g.List {
		comments[i] = string(c.Text)
	}

	lines := make([]string, 0, 10) // most comments are less than 10 lines
	for _, c := range comments {
		// Remove comment markers.
		// The parser has given us exactly the comment text.
		switch c[1] {
		case '/':
			//-style comment (no newline at the end)
			c = c[2:]
			// strip first space - required for Example tests
			if len(c) > 0 && c[0] == ' ' {
				c = c[1:]
			}
		case '*':
			/*-style comment */
			c = c[2 : len(c)-2]
		}

		// Split on newlines.
		cl := strings.Split(c, "\n")

		// Walk lines, stripping trailing white space and adding to list.
		for _, l := range cl {
			lines = append(lines, stripTrailingWhitespace(l))
		}
	}

	// Remove leading blank lines; convert runs of
	// interior blank lines to a single blank line.
	n := 0
	for _, line := range lines {
		if line != "" || n > 0 && lines[n-1] != "" {
			lines[n] = line
			n++
		}
	}
	lines = lines[0:n]

	// Add final "" entry to get trailing newline from Join.
	if n > 0 && lines[n-1] != "" {
		lines = append(lines, "")
	}

	return strings.Join(lines, "\n")
}

// ----------------------------------------------------------------------------
// Expressions and types

// A Field represents a Field declaration list in a struct type,
// a method list in an interface type, or a parameter/result declaration
// in a signature.
//
type Field struct {
	Doc     *CommentGroup // associated documentation; or nil
	Names   []*Ident      // field/method/parameter names; or nil if anonymous field
	Type    Expr          // field/method/parameter type
	Tag     *BasicLit     // field tag; or nil
	Comment *CommentGroup // line comments; or nil
}

func (f *Field) Pos() go_token.Pos {
	if len(f.Names) > 0 {
		return f.Names[0].Pos()
	}
	return f.Type.Pos()
}

func (f *Field) End() go_token.Pos {
	if f.Tag != nil {
		return f.Tag.End()
	}
	return f.Type.End()
}

// A FieldList represents a list of Fields, enclosed by parentheses or braces.
type FieldList struct {
	Opening go_token.Pos // position of opening parenthesis/brace, if any
	List    []*Field     // field list; or nil
	Closing go_token.Pos // position of closing parenthesis/brace, if any
}

func (f *FieldList) Pos() go_token.Pos {
	if f.Opening.IsValid() {
		return f.Opening
	}
	// the list should not be empty in this case;
	// be conservative and guard against bad ASTs
	if len(f.List) > 0 {
		return f.List[0].Pos()
	}
	return go_token.NoPos
}

func (f *FieldList) End() go_token.Pos {
	if f.Closing.IsValid() {
		return f.Closing + 1
	}
	// the list should not be empty in this case;
	// be conservative and guard against bad ASTs
	if n := len(f.List); n > 0 {
		return f.List[n-1].End()
	}
	return go_token.NoPos
}

// NumFields returns the number of (named and anonymous fields) in a FieldList.
func (f *FieldList) NumFields() int {
	n := 0
	if f != nil {
		for _, g := range f.List {
			m := len(g.Names)
			if m == 0 {
				m = 1 // anonymous field
			}
			n += m
		}
	}
	return n
}

// An expression is represented by a tree consisting of one
// or more of the following concrete expression nodes.
//
type (
	// A BadExpr node is a placeholder for expressions containing
	// syntax errors for which no correct expression nodes can be
	// created.
	//
	BadExpr struct {
		From, To go_token.Pos // position range of bad expression
	}

	// An Ident node represents an identifier.
	Ident struct {
		NamePos go_token.Pos   // identifier position
		Name    string         // identifier name
		Obj     *go_ast.Object // denoted object; or nil
	}

	// An Ellipsis node stands for the "..." type in a
	// parameter list or the "..." length in an array type.
	//
	Ellipsis struct {
		Ellipsis go_token.Pos // position of "..."
		Elt      Expr         // ellipsis element type (parameter lists only); or nil
	}

	// A BasicLit node represents a literal of basic type.
	BasicLit struct {
		ValuePos go_token.Pos // literal position
		Kind     token.Token  // token.INT, token.FLOAT, token.IMAG, token.CHAR, or token.STRING
		Value    string       // literal string; e.g. 42, 0x7f, 3.14, 1e-9, 2.4i, 'a', '\x7f', "foo" or `\m\n\o`
	}

	// A FuncLit node represents a function literal.
	FuncLit struct {
		Type *FuncType  // function type
		Body *BlockStmt // function body
	}

	// A CompositeLit node represents a composite literal.
	CompositeLit struct {
		Type   Expr         // literal type; or nil
		Lbrace go_token.Pos // position of "{"
		Elts   []Expr       // list of composite elements; or nil
		Rbrace go_token.Pos // position of "}"
	}

	// A ParenExpr node represents a parenthesized expression.
	ParenExpr struct {
		Lparen go_token.Pos // position of "("
		X      Expr         // parenthesized expression
		Rparen go_token.Pos // position of ")"
	}

	// A SelectorExpr node represents an expression followed by a selector.
	SelectorExpr struct {
		X   Expr   // expression
		Sel *Ident // field selector
	}

	// An IndexExpr node represents an expression followed by an index.
	IndexExpr struct {
		X      Expr         // expression
		Lbrack go_token.Pos // position of "["
		Index  Expr         // index expression
		Rbrack go_token.Pos // position of "]"
	}

	// An SliceExpr node represents an expression followed by slice indices.
	SliceExpr struct {
		X      Expr         // expression
		Lbrack go_token.Pos // position of "["
		Low    Expr         // begin of slice range; or nil
		High   Expr         // end of slice range; or nil
		Rbrack go_token.Pos // position of "]"
	}

	// A TypeAssertExpr node represents an expression followed by a
	// type assertion.
	//
	TypeAssertExpr struct {
		X    Expr // expression
		Type Expr // asserted type; nil means type switch X.(type)
	}

	// A CallExpr node represents an expression followed by an argument list.
	CallExpr struct {
		Fun      Expr         // function expression
		Lparen   go_token.Pos // position of "("
		Args     []Expr       // function arguments; or nil
		Ellipsis go_token.Pos // position of "...", if any
		Rparen   go_token.Pos // position of ")"
	}

	// A StarExpr node represents an expression of the form "*" Expression.
	// Semantically it could be a unary "*" expression, or a pointer type.
	//
	StarExpr struct {
		Star go_token.Pos // position of "*"
		X    Expr         // operand
	}

	// A UnaryExpr node represents a unary expression.
	// Unary "*" expressions are represented via StarExpr nodes.
	//
	UnaryExpr struct {
		OpPos go_token.Pos // position of Op
		Op    token.Token  // operator
		X     Expr         // operand
	}

	// A BinaryExpr node represents a binary expression.
	BinaryExpr struct {
		X     Expr         // left operand
		OpPos go_token.Pos // position of Op
		Op    token.Token  // operator
		Y     Expr         // right operand
	}

	// A KeyValueExpr node represents (key : value) pairs
	// in composite literals.
	//
	KeyValueExpr struct {
		Key   Expr
		Colon go_token.Pos // position of ":"
		Value Expr
	}
)

// The direction of a channel type is indicated by one
// of the following constants.
//
type ChanDir int

const (
	SEND ChanDir = 1 << iota
	RECV
)

// A type is represented by a tree consisting of one
// or more of the following type-specific expression
// nodes.
//
type (
	// An ArrayType node represents an array or slice type.
	ArrayType struct {
		Lbrack go_token.Pos // position of "["
		Len    Expr         // Ellipsis node for [...]T array types, nil for slice types
		Elt    Expr         // element type
	}

	// A StructType node represents a struct type.
	StructType struct {
		Struct     go_token.Pos // position of "struct" keyword
		Fields     *FieldList   // list of field declarations
		Incomplete bool         // true if (source) fields are missing in the Fields list
	}

	// Pointer types are represented via StarExpr nodes.

	// A FuncType node represents a function type.
	FuncType struct {
		Func    go_token.Pos // position of "func" keyword
		Params  *FieldList   // (incoming) parameters; or nil
		Results *FieldList   // (outgoing) results; or nil
	}

	// An InterfaceType node represents an interface type.
	InterfaceType struct {
		Interface  go_token.Pos // position of "interface" keyword
		Methods    *FieldList   // list of methods
		Incomplete bool         // true if (source) methods are missing in the Methods list
	}

	// A MapType node represents a map type.
	MapType struct {
		Map   go_token.Pos // position of "map" keyword
		Key   Expr
		Value Expr
	}

	// A ViewType node represents a view type.
	ViewType struct {
		View  go_token.Pos // position of "view" keyword
		Key   Expr         // primary key type or nil
		Value Expr
	}

	// A ChanType node represents a channel type.
	ChanType struct {
		Begin go_token.Pos // position of "chan" keyword or "<-" (whichever comes first)
		Dir   ChanDir      // channel direction
		Value Expr         // value type
	}
)

type (
	Step interface {
		Expr
		GoStep() *go_ast.CallExpr
		stepNode()
	}

	// type V view[]M
	// source(M) => V
	SourceStep struct {
		Source go_token.Pos // position of the "source" keyword
		Lparen go_token.Pos
		Type   Expr
		Rparen go_token.Pos
	}

	// type V view[...]M
	// V.select(func(M)bool) => V
	SelectStep struct {
		X      Expr
		Select go_token.Pos // position of the "select" keyword
		Lparen go_token.Pos
		F      Expr
		Rparen go_token.Pos
	}

	// type V view[...]M
	// V.reject(func(M)bool) => V
	RejectStep struct {
		X      Expr
		Reject go_token.Pos // position of the "reject" keyword
		Lparen go_token.Pos
		F      Expr
		Rparen go_token.Pos
	}

	// type V view[...]M
	// V.detect(func(M)bool) => V
	DetectStep struct {
		X      Expr
		Detect go_token.Pos // position of the "detect" keyword
		Lparen go_token.Pos
		F      Expr
		Rparen go_token.Pos
	}

	// type V view[...]M
	// type W view[...]N
	// V.collect(func(M)N) => W
	CollectStep struct {
		X       Expr
		Collect go_token.Pos // position of the "collect" keyword
		Lparen  go_token.Pos
		F       Expr
		Rparen  go_token.Pos
	}

	// type V view[...]M
	// V.inject(func(M, []A)A) => A
	InjectStep struct {
		X      Expr
		Inject go_token.Pos // position of the "inject" keyword
		Lparen go_token.Pos
		F      Expr
		Rparen go_token.Pos
	}

	// type V view[...]M
	// type W view[K]view[...]M
	// V.group(func(M)K) => W
	GroupStep struct {
		X      Expr
		Group  go_token.Pos // position of the "group" keyword
		Lparen go_token.Pos
		F      Expr
		Rparen go_token.Pos
	}

	// type V view[...]M
	// type W view[I]M
	// V.index(func(M)I) => W
	IndexStep struct {
		X      Expr
		Index  go_token.Pos // position of the "index" keyword
		Lparen go_token.Pos
		F      Expr
		Rparen go_token.Pos
	}

	// type V view[...]M
	// V.sort(func(M)I) => V
	SortStep struct {
		X      Expr
		Sort   go_token.Pos // position of the "sort" keyword
		Lparen go_token.Pos
		F      Expr
		Rparen go_token.Pos
	}
)

// Pos and End implementations for expression/type nodes.
//
func (x *BadExpr) Pos() go_token.Pos  { return x.From }
func (x *Ident) Pos() go_token.Pos    { return x.NamePos }
func (x *Ellipsis) Pos() go_token.Pos { return x.Ellipsis }
func (x *BasicLit) Pos() go_token.Pos { return x.ValuePos }
func (x *FuncLit) Pos() go_token.Pos  { return x.Type.Pos() }
func (x *CompositeLit) Pos() go_token.Pos {
	if x.Type != nil {
		return x.Type.Pos()
	}
	return x.Lbrace
}
func (x *ParenExpr) Pos() go_token.Pos      { return x.Lparen }
func (x *SelectorExpr) Pos() go_token.Pos   { return x.X.Pos() }
func (x *IndexExpr) Pos() go_token.Pos      { return x.X.Pos() }
func (x *SliceExpr) Pos() go_token.Pos      { return x.X.Pos() }
func (x *TypeAssertExpr) Pos() go_token.Pos { return x.X.Pos() }
func (x *CallExpr) Pos() go_token.Pos       { return x.Fun.Pos() }
func (x *StarExpr) Pos() go_token.Pos       { return x.Star }
func (x *UnaryExpr) Pos() go_token.Pos      { return x.OpPos }
func (x *BinaryExpr) Pos() go_token.Pos     { return x.X.Pos() }
func (x *KeyValueExpr) Pos() go_token.Pos   { return x.Key.Pos() }
func (x *ArrayType) Pos() go_token.Pos      { return x.Lbrack }
func (x *StructType) Pos() go_token.Pos     { return x.Struct }
func (x *FuncType) Pos() go_token.Pos       { return x.Func }
func (x *InterfaceType) Pos() go_token.Pos  { return x.Interface }
func (x *MapType) Pos() go_token.Pos        { return x.Map }
func (x *ViewType) Pos() go_token.Pos       { return x.View }
func (x *ChanType) Pos() go_token.Pos       { return x.Begin }

func (x *BadExpr) End() go_token.Pos { return x.To }
func (x *Ident) End() go_token.Pos   { return go_token.Pos(int(x.NamePos) + len(x.Name)) }
func (x *Ellipsis) End() go_token.Pos {
	if x.Elt != nil {
		return x.Elt.End()
	}
	return x.Ellipsis + 3 // len("...")
}
func (x *BasicLit) End() go_token.Pos     { return go_token.Pos(int(x.ValuePos) + len(x.Value)) }
func (x *FuncLit) End() go_token.Pos      { return x.Body.End() }
func (x *CompositeLit) End() go_token.Pos { return x.Rbrace + 1 }
func (x *ParenExpr) End() go_token.Pos    { return x.Rparen + 1 }
func (x *SelectorExpr) End() go_token.Pos { return x.Sel.End() }
func (x *IndexExpr) End() go_token.Pos    { return x.Rbrack + 1 }
func (x *SliceExpr) End() go_token.Pos    { return x.Rbrack + 1 }
func (x *TypeAssertExpr) End() go_token.Pos {
	if x.Type != nil {
		return x.Type.End()
	}
	return x.X.End()
}
func (x *CallExpr) End() go_token.Pos     { return x.Rparen + 1 }
func (x *StarExpr) End() go_token.Pos     { return x.X.End() }
func (x *UnaryExpr) End() go_token.Pos    { return x.X.End() }
func (x *BinaryExpr) End() go_token.Pos   { return x.Y.End() }
func (x *KeyValueExpr) End() go_token.Pos { return x.Value.End() }
func (x *ArrayType) End() go_token.Pos    { return x.Elt.End() }
func (x *StructType) End() go_token.Pos   { return x.Fields.End() }
func (x *FuncType) End() go_token.Pos {
	if x.Results != nil {
		return x.Results.End()
	}
	return x.Params.End()
}
func (x *InterfaceType) End() go_token.Pos { return x.Methods.End() }
func (x *MapType) End() go_token.Pos       { return x.Value.End() }
func (x *ViewType) End() go_token.Pos      { return x.Value.End() }
func (x *ChanType) End() go_token.Pos      { return x.Value.End() }

// exprNode() ensures that only expression/type nodes can be
// assigned to an ExprNode.
//
func (*BadExpr) exprNode()        {}
func (*Ident) exprNode()          {}
func (*Ellipsis) exprNode()       {}
func (*BasicLit) exprNode()       {}
func (*FuncLit) exprNode()        {}
func (*CompositeLit) exprNode()   {}
func (*ParenExpr) exprNode()      {}
func (*SelectorExpr) exprNode()   {}
func (*IndexExpr) exprNode()      {}
func (*SliceExpr) exprNode()      {}
func (*TypeAssertExpr) exprNode() {}
func (*CallExpr) exprNode()       {}
func (*StarExpr) exprNode()       {}
func (*UnaryExpr) exprNode()      {}
func (*BinaryExpr) exprNode()     {}
func (*KeyValueExpr) exprNode()   {}

func (*ArrayType) exprNode()     {}
func (*StructType) exprNode()    {}
func (*FuncType) exprNode()      {}
func (*InterfaceType) exprNode() {}
func (*MapType) exprNode()       {}
func (*ViewType) exprNode()      {}
func (*ChanType) exprNode()      {}

func (x *SourceStep) Pos() go_token.Pos  { return x.Source }
func (x *SelectStep) Pos() go_token.Pos  { return x.Select }
func (x *RejectStep) Pos() go_token.Pos  { return x.Reject }
func (x *DetectStep) Pos() go_token.Pos  { return x.Detect }
func (x *CollectStep) Pos() go_token.Pos { return x.Collect }
func (x *InjectStep) Pos() go_token.Pos  { return x.Inject }
func (x *GroupStep) Pos() go_token.Pos   { return x.Group }
func (x *IndexStep) Pos() go_token.Pos   { return x.Index }
func (x *SortStep) Pos() go_token.Pos    { return x.Sort }

func (x *SourceStep) End() go_token.Pos  { return x.Rparen }
func (x *SelectStep) End() go_token.Pos  { return x.Rparen }
func (x *RejectStep) End() go_token.Pos  { return x.Rparen }
func (x *DetectStep) End() go_token.Pos  { return x.Rparen }
func (x *CollectStep) End() go_token.Pos { return x.Rparen }
func (x *InjectStep) End() go_token.Pos  { return x.Rparen }
func (x *GroupStep) End() go_token.Pos   { return x.Rparen }
func (x *IndexStep) End() go_token.Pos   { return x.Rparen }
func (x *SortStep) End() go_token.Pos    { return x.Rparen }

func (*SourceStep) exprNode()  {}
func (*SelectStep) exprNode()  {}
func (*RejectStep) exprNode()  {}
func (*DetectStep) exprNode()  {}
func (*CollectStep) exprNode() {}
func (*InjectStep) exprNode()  {}
func (*GroupStep) exprNode()   {}
func (*IndexStep) exprNode()   {}
func (*SortStep) exprNode()    {}

func (*SourceStep) stepNode()  {}
func (*SelectStep) stepNode()  {}
func (*RejectStep) stepNode()  {}
func (*DetectStep) stepNode()  {}
func (*CollectStep) stepNode() {}
func (*InjectStep) stepNode()  {}
func (*GroupStep) stepNode()   {}
func (*IndexStep) stepNode()   {}
func (*SortStep) stepNode()    {}

// ----------------------------------------------------------------------------
// Convenience functions for Idents

var noPos go_token.Pos

// NewIdent creates a new Ident without position.
// Useful for ASTs generated by code other than the Go parser.
//
func NewIdent(name string) *Ident { return &Ident{noPos, name, nil} }

// IsExported returns whether name is an exported Go symbol
// (i.e., whether it begins with an uppercase letter).
//
func IsExported(name string) bool {
	ch, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(ch)
}

// IsExported returns whether id is an exported Go symbol
// (i.e., whether it begins with an uppercase letter).
//
func (id *Ident) IsExported() bool { return IsExported(id.Name) }

func (id *Ident) String() string {
	if id != nil {
		return id.Name
	}
	return "<nil>"
}

// ----------------------------------------------------------------------------
// Statements

// A statement is represented by a tree consisting of one
// or more of the following concrete statement nodes.
//
type (
	// A BadStmt node is a placeholder for statements containing
	// syntax errors for which no correct statement nodes can be
	// created.
	//
	BadStmt struct {
		From, To go_token.Pos // position range of bad statement
	}

	// A DeclStmt node represents a declaration in a statement list.
	DeclStmt struct {
		Decl Decl
	}

	// An EmptyStmt node represents an empty statement.
	// The "position" of the empty statement is the position
	// of the immediately preceding semicolon.
	//
	EmptyStmt struct {
		Semicolon go_token.Pos // position of preceding ";"
	}

	// A LabeledStmt node represents a labeled statement.
	LabeledStmt struct {
		Label *Ident
		Colon go_token.Pos // position of ":"
		Stmt  Stmt
	}

	// An ExprStmt node represents a (stand-alone) expression
	// in a statement list.
	//
	ExprStmt struct {
		X Expr // expression
	}

	// A SendStmt node represents a send statement.
	SendStmt struct {
		Chan  Expr
		Arrow go_token.Pos // position of "<-"
		Value Expr
	}

	// An IncDecStmt node represents an increment or decrement statement.
	IncDecStmt struct {
		X      Expr
		TokPos go_token.Pos // position of Tok
		Tok    token.Token  // INC or DEC
	}

	// An AssignStmt node represents an assignment or
	// a short variable declaration.
	//
	AssignStmt struct {
		Lhs    []Expr
		TokPos go_token.Pos // position of Tok
		Tok    token.Token  // assignment token, DEFINE
		Rhs    []Expr
	}

	// A GoStmt node represents a go statement.
	GoStmt struct {
		Go   go_token.Pos // position of "go" keyword
		Call *CallExpr
	}

	// A DeferStmt node represents a defer statement.
	DeferStmt struct {
		Defer go_token.Pos // position of "defer" keyword
		Call  *CallExpr
	}

	// A ReturnStmt node represents a return statement.
	ReturnStmt struct {
		Return  go_token.Pos // position of "return" keyword
		Results []Expr       // result expressions; or nil
	}

	// A BranchStmt node represents a break, continue, goto,
	// or fallthrough statement.
	//
	BranchStmt struct {
		TokPos go_token.Pos // position of Tok
		Tok    token.Token  // keyword token (BREAK, CONTINUE, GOTO, FALLTHROUGH)
		Label  *Ident       // label name; or nil
	}

	// A BlockStmt node represents a braced statement list.
	BlockStmt struct {
		Lbrace go_token.Pos // position of "{"
		List   []Stmt
		Rbrace go_token.Pos // position of "}"
	}

	// An IfStmt node represents an if statement.
	IfStmt struct {
		If   go_token.Pos // position of "if" keyword
		Init Stmt         // initialization statement; or nil
		Cond Expr         // condition
		Body *BlockStmt
		Else Stmt // else branch; or nil
	}

	// A CaseClause represents a case of an expression or type switch statement.
	CaseClause struct {
		Case  go_token.Pos // position of "case" or "default" keyword
		List  []Expr       // list of expressions or types; nil means default case
		Colon go_token.Pos // position of ":"
		Body  []Stmt       // statement list; or nil
	}

	// A SwitchStmt node represents an expression switch statement.
	SwitchStmt struct {
		Switch go_token.Pos // position of "switch" keyword
		Init   Stmt         // initialization statement; or nil
		Tag    Expr         // tag expression; or nil
		Body   *BlockStmt   // CaseClauses only
	}

	// An TypeSwitchStmt node represents a type switch statement.
	TypeSwitchStmt struct {
		Switch go_token.Pos // position of "switch" keyword
		Init   Stmt         // initialization statement; or nil
		Assign Stmt         // x := y.(type) or y.(type)
		Body   *BlockStmt   // CaseClauses only
	}

	// A CommClause node represents a case of a select statement.
	CommClause struct {
		Case  go_token.Pos // position of "case" or "default" keyword
		Comm  Stmt         // send or receive statement; nil means default case
		Colon go_token.Pos // position of ":"
		Body  []Stmt       // statement list; or nil
	}

	// An SelectStmt node represents a select statement.
	SelectStmt struct {
		Select go_token.Pos // position of "select" keyword
		Body   *BlockStmt   // CommClauses only
	}

	// A ForStmt represents a for statement.
	ForStmt struct {
		For  go_token.Pos // position of "for" keyword
		Init Stmt         // initialization statement; or nil
		Cond Expr         // condition; or nil
		Post Stmt         // post iteration statement; or nil
		Body *BlockStmt
	}

	// A RangeStmt represents a for statement with a range clause.
	RangeStmt struct {
		For        go_token.Pos // position of "for" keyword
		Key, Value Expr         // Value may be nil
		TokPos     go_token.Pos // position of Tok
		Tok        token.Token  // ASSIGN, DEFINE
		X          Expr         // value to range over
		Body       *BlockStmt
	}
)

// Pos and End implementations for statement nodes.
//
func (s *BadStmt) Pos() go_token.Pos        { return s.From }
func (s *DeclStmt) Pos() go_token.Pos       { return s.Decl.Pos() }
func (s *EmptyStmt) Pos() go_token.Pos      { return s.Semicolon }
func (s *LabeledStmt) Pos() go_token.Pos    { return s.Label.Pos() }
func (s *ExprStmt) Pos() go_token.Pos       { return s.X.Pos() }
func (s *SendStmt) Pos() go_token.Pos       { return s.Chan.Pos() }
func (s *IncDecStmt) Pos() go_token.Pos     { return s.X.Pos() }
func (s *AssignStmt) Pos() go_token.Pos     { return s.Lhs[0].Pos() }
func (s *GoStmt) Pos() go_token.Pos         { return s.Go }
func (s *DeferStmt) Pos() go_token.Pos      { return s.Defer }
func (s *ReturnStmt) Pos() go_token.Pos     { return s.Return }
func (s *BranchStmt) Pos() go_token.Pos     { return s.TokPos }
func (s *BlockStmt) Pos() go_token.Pos      { return s.Lbrace }
func (s *IfStmt) Pos() go_token.Pos         { return s.If }
func (s *CaseClause) Pos() go_token.Pos     { return s.Case }
func (s *SwitchStmt) Pos() go_token.Pos     { return s.Switch }
func (s *TypeSwitchStmt) Pos() go_token.Pos { return s.Switch }
func (s *CommClause) Pos() go_token.Pos     { return s.Case }
func (s *SelectStmt) Pos() go_token.Pos     { return s.Select }
func (s *ForStmt) Pos() go_token.Pos        { return s.For }
func (s *RangeStmt) Pos() go_token.Pos      { return s.For }

func (s *BadStmt) End() go_token.Pos  { return s.To }
func (s *DeclStmt) End() go_token.Pos { return s.Decl.End() }
func (s *EmptyStmt) End() go_token.Pos {
	return s.Semicolon + 1 /* len(";") */
}
func (s *LabeledStmt) End() go_token.Pos { return s.Stmt.End() }
func (s *ExprStmt) End() go_token.Pos    { return s.X.End() }
func (s *SendStmt) End() go_token.Pos    { return s.Value.End() }
func (s *IncDecStmt) End() go_token.Pos {
	return s.TokPos + 2 /* len("++") */
}
func (s *AssignStmt) End() go_token.Pos { return s.Rhs[len(s.Rhs)-1].End() }
func (s *GoStmt) End() go_token.Pos     { return s.Call.End() }
func (s *DeferStmt) End() go_token.Pos  { return s.Call.End() }
func (s *ReturnStmt) End() go_token.Pos {
	if n := len(s.Results); n > 0 {
		return s.Results[n-1].End()
	}
	return s.Return + 6 // len("return")
}
func (s *BranchStmt) End() go_token.Pos {
	if s.Label != nil {
		return s.Label.End()
	}
	return go_token.Pos(int(s.TokPos) + len(s.Tok.String()))
}
func (s *BlockStmt) End() go_token.Pos { return s.Rbrace + 1 }
func (s *IfStmt) End() go_token.Pos {
	if s.Else != nil {
		return s.Else.End()
	}
	return s.Body.End()
}
func (s *CaseClause) End() go_token.Pos {
	if n := len(s.Body); n > 0 {
		return s.Body[n-1].End()
	}
	return s.Colon + 1
}
func (s *SwitchStmt) End() go_token.Pos     { return s.Body.End() }
func (s *TypeSwitchStmt) End() go_token.Pos { return s.Body.End() }
func (s *CommClause) End() go_token.Pos {
	if n := len(s.Body); n > 0 {
		return s.Body[n-1].End()
	}
	return s.Colon + 1
}
func (s *SelectStmt) End() go_token.Pos { return s.Body.End() }
func (s *ForStmt) End() go_token.Pos    { return s.Body.End() }
func (s *RangeStmt) End() go_token.Pos  { return s.Body.End() }

// stmtNode() ensures that only statement nodes can be
// assigned to a StmtNode.
//
func (*BadStmt) stmtNode()        {}
func (*DeclStmt) stmtNode()       {}
func (*EmptyStmt) stmtNode()      {}
func (*LabeledStmt) stmtNode()    {}
func (*ExprStmt) stmtNode()       {}
func (*SendStmt) stmtNode()       {}
func (*IncDecStmt) stmtNode()     {}
func (*AssignStmt) stmtNode()     {}
func (*GoStmt) stmtNode()         {}
func (*DeferStmt) stmtNode()      {}
func (*ReturnStmt) stmtNode()     {}
func (*BranchStmt) stmtNode()     {}
func (*BlockStmt) stmtNode()      {}
func (*IfStmt) stmtNode()         {}
func (*CaseClause) stmtNode()     {}
func (*SwitchStmt) stmtNode()     {}
func (*TypeSwitchStmt) stmtNode() {}
func (*CommClause) stmtNode()     {}
func (*SelectStmt) stmtNode()     {}
func (*ForStmt) stmtNode()        {}
func (*RangeStmt) stmtNode()      {}

// ----------------------------------------------------------------------------
// Declarations

// A Spec node represents a single (non-parenthesized) import,
// constant, type, or variable declaration.
//
type (
	// The Spec type stands for any of *ImportSpec, *ValueSpec, and *TypeSpec.
	Spec interface {
		Node
		GoSpec() go_ast.Spec
		specNode()
	}

	// An ImportSpec node represents a single package import.
	ImportSpec struct {
		Doc     *CommentGroup // associated documentation; or nil
		Name    *Ident        // local package name (including "."); or nil
		Path    *BasicLit     // import path
		Comment *CommentGroup // line comments; or nil
		EndPos  go_token.Pos  // end of spec (overrides Path.Pos if nonzero)
	}

	// A ValueSpec node represents a constant or variable declaration
	// (ConstSpec or VarSpec production).
	//
	ValueSpec struct {
		Doc     *CommentGroup // associated documentation; or nil
		Names   []*Ident      // value names (len(Names) > 0)
		Type    Expr          // value type; or nil
		Values  []Expr        // initial values; or nil
		Comment *CommentGroup // line comments; or nil
	}

	// A TypeSpec node represents a type declaration (TypeSpec production).
	TypeSpec struct {
		Doc     *CommentGroup // associated documentation; or nil
		Name    *Ident        // type name
		Type    Expr          // *Ident, *ParenExpr, *SelectorExpr, *StarExpr, or any of the *XxxTypes
		Comment *CommentGroup // line comments; or nil
	}
)

// Pos and End implementations for spec nodes.
//
func (s *ImportSpec) Pos() go_token.Pos {
	if s.Name != nil {
		return s.Name.Pos()
	}
	return s.Path.Pos()
}
func (s *ValueSpec) Pos() go_token.Pos { return s.Names[0].Pos() }
func (s *TypeSpec) Pos() go_token.Pos  { return s.Name.Pos() }

func (s *ImportSpec) End() go_token.Pos {
	if s.EndPos != 0 {
		return s.EndPos
	}
	return s.Path.End()
}

func (s *ValueSpec) End() go_token.Pos {
	if n := len(s.Values); n > 0 {
		return s.Values[n-1].End()
	}
	if s.Type != nil {
		return s.Type.End()
	}
	return s.Names[len(s.Names)-1].End()
}
func (s *TypeSpec) End() go_token.Pos { return s.Type.End() }

// specNode() ensures that only spec nodes can be
// assigned to a Spec.
//
func (*ImportSpec) specNode() {}
func (*ValueSpec) specNode()  {}
func (*TypeSpec) specNode()   {}

// A declaration is represented by one of the following declaration nodes.
//
type (
	// A BadDecl node is a placeholder for declarations containing
	// syntax errors for which no correct declaration nodes can be
	// created.
	//
	BadDecl struct {
		From, To go_token.Pos // position range of bad declaration
	}

	// A GenDecl node (generic declaration node) represents an import,
	// constant, type or variable declaration. A valid Lparen position
	// (Lparen.Line > 0) indicates a parenthesized declaration.
	//
	// Relationship between Tok value and Specs element type:
	//
	//	token.IMPORT  *ImportSpec
	//	token.CONST   *ValueSpec
	//	token.TYPE    *TypeSpec
	//	token.VAR     *ValueSpec
	//
	GenDecl struct {
		Doc    *CommentGroup // associated documentation; or nil
		TokPos go_token.Pos  // position of Tok
		Tok    token.Token   // IMPORT, CONST, TYPE, VAR
		Lparen go_token.Pos  // position of '(', if any
		Specs  []Spec
		Rparen go_token.Pos // position of ')', if any
	}

	// A FuncDecl node represents a function declaration.
	FuncDecl struct {
		Doc  *CommentGroup // associated documentation; or nil
		Recv *FieldList    // receiver (methods); or nil (functions)
		Name *Ident        // function/method name
		Type *FuncType     // position of Func keyword, parameters and results
		Body *BlockStmt    // function body; or nil (forward declaration)
	}
)

// Pos and End implementations for declaration nodes.
//
func (d *BadDecl) Pos() go_token.Pos  { return d.From }
func (d *GenDecl) Pos() go_token.Pos  { return d.TokPos }
func (d *FuncDecl) Pos() go_token.Pos { return d.Type.Pos() }

func (d *BadDecl) End() go_token.Pos { return d.To }
func (d *GenDecl) End() go_token.Pos {
	if d.Rparen.IsValid() {
		return d.Rparen + 1
	}
	return d.Specs[0].End()
}
func (d *FuncDecl) End() go_token.Pos {
	if d.Body != nil {
		return d.Body.End()
	}
	return d.Type.End()
}

// declNode() ensures that only declaration nodes can be
// assigned to a DeclNode.
//
func (*BadDecl) declNode()  {}
func (*GenDecl) declNode()  {}
func (*FuncDecl) declNode() {}

// ----------------------------------------------------------------------------
// Files and packages

// A File node represents a Go source file.
//
// The Comments list contains all comments in the source file in order of
// appearance, including the comments that are pointed to from other nodes
// via Doc and Comment fields.
//
type File struct {
	Doc        *CommentGroup   // associated documentation; or nil
	Package    go_token.Pos    // position of "package" keyword
	Name       *Ident          // package name
	Decls      []Decl          // top-level declarations; or nil
	Scope      *go_ast.Scope   // package scope (this file only)
	Imports    []*ImportSpec   // imports in this file
	Unresolved []*Ident        // unresolved identifiers in this file
	Comments   []*CommentGroup // list of all comments in the source file
}

func (f *File) Pos() go_token.Pos { return f.Package }
func (f *File) End() go_token.Pos {
	if n := len(f.Decls); n > 0 {
		return f.Decls[n-1].End()
	}
	return f.Name.End()
}

// A Package node represents a set of source files
// collectively building a Go package.
//
type Package struct {
	Name    string                    // package name
	Scope   *go_ast.Scope             // package scope across all files
	Imports map[string]*go_ast.Object // map of package id -> package object
	Files   map[string]*File          // Go source files by filename
}

func (p *Package) Pos() go_token.Pos { return go_token.NoPos }
func (p *Package) End() go_token.Pos { return go_token.NoPos }

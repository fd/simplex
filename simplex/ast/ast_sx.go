package ast

import (
	"github.com/fd/w/simplex/token"
	//go_ast "go/ast"
)

type (

	// A ViewType node represents a view type.
	ViewType struct {
		View  token.Pos // position of "view" keyword
		Key   Expr      // primary key type or nil
		Value Expr
	}
)

type (
	Step interface {
		Expr
		//GoStep() *go_ast.CallExpr
		stepNode()
	}

	// type V view[]M
	// source(M) => V
	SourceStep struct {
		Source token.Pos // position of the "source" keyword
		Lparen token.Pos
		Type   Expr
		Rparen token.Pos
	}

	// type V view[...]M
	// V.select(func(M)bool) => V
	SelectStep struct {
		X      Expr
		Select token.Pos // position of the "select" keyword
		Lparen token.Pos
		F      Expr
		Rparen token.Pos
	}

	// type V view[...]M
	// V.reject(func(M)bool) => V
	RejectStep struct {
		X      Expr
		Reject token.Pos // position of the "reject" keyword
		Lparen token.Pos
		F      Expr
		Rparen token.Pos
	}

	// type V view[...]M
	// V.detect(func(M)bool) => V
	DetectStep struct {
		X      Expr
		Detect token.Pos // position of the "detect" keyword
		Lparen token.Pos
		F      Expr
		Rparen token.Pos
	}
	// type V view[...]M
	// type W view[...]N
	// V.collect(func(M)N) => W
	CollectStep struct {
		X       Expr
		Collect token.Pos // position of the "collect" keyword
		Lparen  token.Pos
		F       Expr
		Rparen  token.Pos
	}

	// type V view[...]M
	// V.inject(func(M, []A)A) => A
	InjectStep struct {
		X      Expr
		Inject token.Pos // position of the "inject" keyword
		Lparen token.Pos
		F      Expr
		Rparen token.Pos
	}

	// type V view[...]M
	// type W view[K]view[...]M
	// V.group(func(M)K) => W
	GroupStep struct {
		X      Expr
		Group  token.Pos // position of the "group" keyword
		Lparen token.Pos
		F      Expr
		Rparen token.Pos
	}

	// type V view[...]M
	// type W view[I]M
	// V.index(func(M)I) => W
	IndexStep struct {
		X      Expr
		Index  token.Pos // position of the "index" keyword
		Lparen token.Pos
		F      Expr
		Rparen token.Pos
	}
	// type V view[...]M
	// V.sort(func(M)I) => V
	SortStep struct {
		X      Expr
		Sort   token.Pos // position of the "sort" keyword
		Lparen token.Pos
		F      Expr
		Rparen token.Pos
	}
)

func (x *ViewType) Pos() token.Pos { return x.View }
func (x *ViewType) End() token.Pos { return x.Value.End() }
func (*ViewType) exprNode()        {}

func (x *SourceStep) Pos() token.Pos  { return x.Source }
func (x *SelectStep) Pos() token.Pos  { return x.Select }
func (x *RejectStep) Pos() token.Pos  { return x.Reject }
func (x *DetectStep) Pos() token.Pos  { return x.Detect }
func (x *CollectStep) Pos() token.Pos { return x.Collect }
func (x *InjectStep) Pos() token.Pos  { return x.Inject }
func (x *GroupStep) Pos() token.Pos   { return x.Group }
func (x *IndexStep) Pos() token.Pos   { return x.Index }
func (x *SortStep) Pos() token.Pos    { return x.Sort }

func (x *SourceStep) End() token.Pos  { return x.Rparen }
func (x *SelectStep) End() token.Pos  { return x.Rparen }
func (x *RejectStep) End() token.Pos  { return x.Rparen }
func (x *DetectStep) End() token.Pos  { return x.Rparen }
func (x *CollectStep) End() token.Pos { return x.Rparen }
func (x *InjectStep) End() token.Pos  { return x.Rparen }
func (x *GroupStep) End() token.Pos   { return x.Rparen }
func (x *IndexStep) End() token.Pos   { return x.Rparen }
func (x *SortStep) End() token.Pos    { return x.Rparen }

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

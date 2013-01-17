package ast

import (
	"github.com/fd/w/simplex/token"
)

type (

	// A ViewType node represents a view type.
	ViewType struct {
		View  token.Pos // position of "view" keyword
		Key   Expr      // primary key type or nil
		Value Expr
	}

	// A TableType node represents a table type.
	TableType struct {
		Table token.Pos // position of "table" keyword
		Key   Expr      // primary key type
		Value Expr
	}

	StepType int
)

func (x *ViewType) Pos() token.Pos { return x.View }
func (x *ViewType) End() token.Pos { return x.Value.End() }
func (*ViewType) exprNode()        {}

func (x *TableType) Pos() token.Pos { return x.Table }
func (x *TableType) End() token.Pos { return x.Value.End() }
func (*TableType) exprNode()        {}

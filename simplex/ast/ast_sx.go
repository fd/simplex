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

const (
	BadStep StepType = iota
	SelectStep
	RejectStep
	DetectStep
	CollectStep
	InjectStep
	GroupStep
	IndexStep
	SortStep
)

var StepTypeNames = map[string]StepType{
	"select":  SelectStep,
	"reject":  RejectStep,
	"detect":  DetectStep,
	"collect": CollectStep,
	"inject":  InjectStep,
	"group":   GroupStep,
	"index":   IndexStep,
	"sort":    SortStep,
}

type (
	StepExpr struct {
		X        Expr
		TokPos   token.Pos // position of the step name keyword
		StepType StepType
		Lparen   token.Pos
		F        Expr
		Rparen   token.Pos
	}
)

func (x *ViewType) Pos() token.Pos { return x.View }
func (x *ViewType) End() token.Pos { return x.Value.End() }
func (*ViewType) exprNode()        {}

func (x *TableType) Pos() token.Pos { return x.Table }
func (x *TableType) End() token.Pos { return x.Value.End() }
func (*TableType) exprNode()        {}

func (x *StepExpr) Pos() token.Pos { return x.TokPos }
func (x *StepExpr) End() token.Pos { return x.Rparen }
func (*StepExpr) exprNode()        {}

func (typ StepType) String() string {
	for name, t := range StepTypeNames {
		if t == typ {
			return name
		}
	}
	return "(bad step)"
}

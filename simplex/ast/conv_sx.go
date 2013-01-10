package ast

import (
	go_ast "go/ast"
)

func (x *SourceStep) GoSourceStep() *go_ast.CallExpr   { panic("not implemented") }
func (x *SelectStep) GoSelectStep() *go_ast.CallExpr   { panic("not implemented") }
func (x *RejectStep) GoRejectStep() *go_ast.CallExpr   { panic("not implemented") }
func (x *DetectStep) GoDetectStep() *go_ast.CallExpr   { panic("not implemented") }
func (x *CollectStep) GoCollectStep() *go_ast.CallExpr { panic("not implemented") }
func (x *InjectStep) GoInjectStep() *go_ast.CallExpr   { panic("not implemented") }
func (x *GroupStep) GoGroupStep() *go_ast.CallExpr     { panic("not implemented") }
func (x *IndexStep) GoIndexStep() *go_ast.CallExpr     { panic("not implemented") }
func (x *SortStep) GoSortStep() *go_ast.CallExpr       { panic("not implemented") }

func (x *SourceStep) GoStep() *go_ast.CallExpr  { return x.GoSourceStep() }
func (x *SelectStep) GoStep() *go_ast.CallExpr  { return x.GoSelectStep() }
func (x *RejectStep) GoStep() *go_ast.CallExpr  { return x.GoRejectStep() }
func (x *DetectStep) GoStep() *go_ast.CallExpr  { return x.GoDetectStep() }
func (x *CollectStep) GoStep() *go_ast.CallExpr { return x.GoCollectStep() }
func (x *InjectStep) GoStep() *go_ast.CallExpr  { return x.GoInjectStep() }
func (x *GroupStep) GoStep() *go_ast.CallExpr   { return x.GoGroupStep() }
func (x *IndexStep) GoStep() *go_ast.CallExpr   { return x.GoIndexStep() }
func (x *SortStep) GoStep() *go_ast.CallExpr    { return x.GoSortStep() }

func (x *SourceStep) GoExpr() go_ast.Expr  { return x.GoStep() }
func (x *SelectStep) GoExpr() go_ast.Expr  { return x.GoStep() }
func (x *RejectStep) GoExpr() go_ast.Expr  { return x.GoStep() }
func (x *DetectStep) GoExpr() go_ast.Expr  { return x.GoStep() }
func (x *CollectStep) GoExpr() go_ast.Expr { return x.GoStep() }
func (x *InjectStep) GoExpr() go_ast.Expr  { return x.GoStep() }
func (x *GroupStep) GoExpr() go_ast.Expr   { return x.GoStep() }
func (x *IndexStep) GoExpr() go_ast.Expr   { return x.GoStep() }
func (x *SortStep) GoExpr() go_ast.Expr    { return x.GoStep() }

func (x *SourceStep) GoNode() go_ast.Node  { return x.GoExpr() }
func (x *SelectStep) GoNode() go_ast.Node  { return x.GoExpr() }
func (x *RejectStep) GoNode() go_ast.Node  { return x.GoExpr() }
func (x *DetectStep) GoNode() go_ast.Node  { return x.GoExpr() }
func (x *CollectStep) GoNode() go_ast.Node { return x.GoExpr() }
func (x *InjectStep) GoNode() go_ast.Node  { return x.GoExpr() }
func (x *GroupStep) GoNode() go_ast.Node   { return x.GoExpr() }
func (x *IndexStep) GoNode() go_ast.Node   { return x.GoExpr() }
func (x *SortStep) GoNode() go_ast.Node    { return x.GoExpr() }

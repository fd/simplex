package parser

import (
	"github.com/fd/w/simplex/ast"
	"github.com/fd/w/simplex/token"
)

func (p *parser) parseViewType() *ast.ViewType {
	if p.trace {
		defer un(trace(p, "ViewType"))
	}

	pos := p.expect(token.VIEW)
	p.expect(token.LBRACK)
	key := p.tryType()
	p.expect(token.RBRACK)
	value := p.parseType()

	return &ast.ViewType{View: pos, Key: key, Value: value}
}

func (p *parser) parseSelectStep(x ast.Expr) ast.Step {
	if p.trace {
		defer un(trace(p, "SelectStep"))
	}

	tok_pos := p.expect(token.SELECT)
	lparen := p.expect(token.LPAREN)

	p.exprLev++
	f := p.parseRhs()
	p.exprLev--

	rparen := p.expectClosing(token.RPAREN, "argument list")

	return &ast.SelectStep{
		X:      x,
		Select: tok_pos,
		Lparen: lparen,
		F:      f,
		Rparen: rparen,
	}
}

func (p *parser) parseRejectStep(x ast.Expr) ast.Step {
	if p.trace {
		defer un(trace(p, "RejectStep"))
	}

	tok_pos := p.expect(token.REJECT)
	lparen := p.expect(token.LPAREN)

	p.exprLev++
	f := p.parseRhs()
	p.exprLev--

	rparen := p.expectClosing(token.RPAREN, "argument list")

	return &ast.RejectStep{
		X:      x,
		Reject: tok_pos,
		Lparen: lparen,
		F:      f,
		Rparen: rparen,
	}
}

func (p *parser) parseDetectStep(x ast.Expr) ast.Step {
	if p.trace {
		defer un(trace(p, "DetectStep"))
	}

	tok_pos := p.expect(token.DETECT)
	lparen := p.expect(token.LPAREN)

	p.exprLev++
	f := p.parseRhs()
	p.exprLev--

	rparen := p.expectClosing(token.RPAREN, "argument list")

	return &ast.DetectStep{
		X:      x,
		Detect: tok_pos,
		Lparen: lparen,
		F:      f,
		Rparen: rparen,
	}
}

func (p *parser) parseCollectStep(x ast.Expr) ast.Step {
	if p.trace {
		defer un(trace(p, "CollectStep"))
	}

	tok_pos := p.expect(token.COLLECT)
	lparen := p.expect(token.LPAREN)

	p.exprLev++
	f := p.parseRhs()
	p.exprLev--

	rparen := p.expectClosing(token.RPAREN, "argument list")

	return &ast.CollectStep{
		X:       x,
		Collect: tok_pos,
		Lparen:  lparen,
		F:       f,
		Rparen:  rparen,
	}
}

func (p *parser) parseInjectStep(x ast.Expr) ast.Step {
	if p.trace {
		defer un(trace(p, "InjectStep"))
	}

	tok_pos := p.expect(token.INJECT)
	lparen := p.expect(token.LPAREN)

	p.exprLev++
	f := p.parseRhs()
	p.exprLev--

	rparen := p.expectClosing(token.RPAREN, "argument list")

	return &ast.InjectStep{
		X:      x,
		Inject: tok_pos,
		Lparen: lparen,
		F:      f,
		Rparen: rparen,
	}
}

func (p *parser) parseGroupStep(x ast.Expr) ast.Step {
	if p.trace {
		defer un(trace(p, "GroupStep"))
	}

	tok_pos := p.expect(token.GROUP)
	lparen := p.expect(token.LPAREN)

	p.exprLev++
	f := p.parseRhs()
	p.exprLev--

	rparen := p.expectClosing(token.RPAREN, "argument list")

	return &ast.GroupStep{
		X:      x,
		Group:  tok_pos,
		Lparen: lparen,
		F:      f,
		Rparen: rparen,
	}
}

func (p *parser) parseIndexStep(x ast.Expr) ast.Step {
	if p.trace {
		defer un(trace(p, "IndexStep"))
	}

	tok_pos := p.expect(token.INDEX)
	lparen := p.expect(token.LPAREN)

	p.exprLev++
	f := p.parseRhs()
	p.exprLev--

	rparen := p.expectClosing(token.RPAREN, "argument list")

	return &ast.IndexStep{
		X:      x,
		Index:  tok_pos,
		Lparen: lparen,
		F:      f,
		Rparen: rparen,
	}
}

func (p *parser) parseSortStep(x ast.Expr) ast.Step {
	if p.trace {
		defer un(trace(p, "SortStep"))
	}

	tok_pos := p.expect(token.SORT)
	lparen := p.expect(token.LPAREN)

	p.exprLev++
	f := p.parseRhs()
	p.exprLev--

	rparen := p.expectClosing(token.RPAREN, "argument list")

	return &ast.SortStep{
		X:      x,
		Sort:   tok_pos,
		Lparen: lparen,
		F:      f,
		Rparen: rparen,
	}
}

// If the result is an identifier, it is not resolved.
//
// NOTE(fd) see parser.go:959
func (p *parser) tryIdentOrType() ast.Expr {
	switch p.tok {
	case token.IDENT:
		return p.parseTypeName()
	case token.LBRACK:
		return p.parseArrayType()
	case token.STRUCT:
		return p.parseStructType()
	case token.MUL:
		return p.parsePointerType()
	case token.FUNC:
		typ, _ := p.parseFuncType()
		return typ
	case token.INTERFACE:
		return p.parseInterfaceType()
	case token.MAP:
		return p.parseMapType()
	case token.VIEW:
		return p.parseViewType()
	case token.CHAN, token.ARROW:
		return p.parseChanType()
	case token.LPAREN:
		lparen := p.pos
		p.next()
		typ := p.parseType()
		rparen := p.expect(token.RPAREN)
		return &ast.ParenExpr{Lparen: lparen, X: typ, Rparen: rparen}
	}

	// no type found
	return nil
}

// parseOperand may return an expression or a raw type (incl. array
// types of the form [...]T. Callers must verify the result.
// If lhs is set and the result is an identifier, it is not resolved.
//
// NOTE(fd) see parser.go:1070
func (p *parser) parseOperand(lhs bool) ast.Expr {
	if p.trace {
		defer un(trace(p, "Operand"))
	}

	switch p.tok {
	case token.IDENT:
		x := p.parseIdent()
		if !lhs {
			p.resolve(x)
		}
		return x

	case token.INT, token.FLOAT, token.IMAG, token.CHAR, token.STRING:
		x := &ast.BasicLit{ValuePos: p.pos, Kind: p.tok, Value: p.lit}
		p.next()
		return x

	case token.LPAREN:
		lparen := p.pos
		p.next()
		p.exprLev++
		x := p.parseRhsOrType() // types may be parenthesized: (some type)
		p.exprLev--
		rparen := p.expect(token.RPAREN)
		return &ast.ParenExpr{Lparen: lparen, X: x, Rparen: rparen}

	case token.FUNC:
		return p.parseFuncTypeOrLit()

	case token.SOURCE:
		pos := p.pos
		p.next()
		lparen := p.expect(token.LPAREN)
		typ := p.parseType()
		rparen := p.expectClosing(token.RPAREN, "argument list")
		return &ast.SourceStep{Source: pos, Lparen: lparen, Type: typ, Rparen: rparen}

	}

	if typ := p.tryIdentOrType(); typ != nil {
		// could be type for composite literal or conversion
		_, isIdent := typ.(*ast.Ident)
		assert(!isIdent, "type cannot be identifier")
		return typ
	}

	// we have an error
	pos := p.pos
	p.errorExpected(pos, "operand")
	syncStmt(p)
	return &ast.BadExpr{From: pos, To: p.pos}
}

// checkExpr checks that x is an expression (and not a type).
//
// NOTE(fd) see parser.go:1275
func (p *parser) checkExpr(x ast.Expr) ast.Expr {
	switch unparen(x).(type) {
	case *ast.BadExpr:
	case *ast.Ident:
	case *ast.BasicLit:
	case *ast.FuncLit:
	case *ast.CompositeLit:
	case *ast.ParenExpr:
		panic("unreachable")
	case *ast.SelectorExpr:
	case *ast.IndexExpr:
	case *ast.SliceExpr:
	case *ast.TypeAssertExpr:
		// If t.Type == nil we have a type assertion of the form
		// y.(type), which is only allowed in type switch expressions.
		// It's hard to exclude those but for the case where we are in
		// a type switch. Instead be lenient and test this in the type
		// checker.
	case *ast.CallExpr:
	case *ast.StarExpr:
	case *ast.UnaryExpr:
	case *ast.BinaryExpr:

	case *ast.SourceStep:
	case *ast.SelectStep:
	case *ast.RejectStep:
	case *ast.DetectStep:
	case *ast.CollectStep:
	case *ast.InjectStep:
	case *ast.GroupStep:
	case *ast.IndexStep:
	case *ast.SortStep:

	default:
		// all other nodes are not proper expressions
		p.errorExpected(x.Pos(), "expression")
		x = &ast.BadExpr{From: x.Pos(), To: x.End()}
	}
	return x
}

// isLiteralType returns true iff x is a legal composite literal type.
//
// NOTE(fd) see parser.go:1322
func isLiteralType(x ast.Expr) bool {
	switch t := x.(type) {
	case *ast.BadExpr:
	case *ast.Ident:
	case *ast.SelectorExpr:
		_, isIdent := t.X.(*ast.Ident)
		return isIdent
	case *ast.ArrayType:
	case *ast.StructType:
	case *ast.MapType:

	case *ast.ViewType:

	default:
		return false // all other nodes are not legal composite literal types
	}
	return true
}

// If lhs is set and the result is an identifier, it is not resolved.
//
// NOTE(fd) see parser.go:1376
func (p *parser) parsePrimaryExpr(lhs bool) ast.Expr {
	if p.trace {
		defer un(trace(p, "PrimaryExpr"))
	}

	x := p.parseOperand(lhs)
L:
	for {
		switch p.tok {
		case token.PERIOD:
			p.next()
			if lhs {
				p.resolve(x)
			}
			switch p.tok {
			case token.IDENT:
				x = p.parseSelector(p.checkExpr(x))

			case token.SELECT:
				x = p.parseSelectStep(p.checkExpr(x))
			case token.REJECT:
				x = p.parseRejectStep(p.checkExpr(x))
			case token.DETECT:
				x = p.parseDetectStep(p.checkExpr(x))
			case token.COLLECT:
				x = p.parseCollectStep(p.checkExpr(x))
			case token.INJECT:
				x = p.parseInjectStep(p.checkExpr(x))
			case token.GROUP:
				x = p.parseGroupStep(p.checkExpr(x))
			case token.INDEX:
				x = p.parseIndexStep(p.checkExpr(x))
			case token.SORT:
				x = p.parseSortStep(p.checkExpr(x))

			case token.LPAREN:
				x = p.parseTypeAssertion(p.checkExpr(x))
			default:
				pos := p.pos
				p.errorExpected(pos, "selector or type assertion")
				p.next() // make progress
				x = &ast.BadExpr{From: pos, To: p.pos}
			}
		case token.LBRACK:
			if lhs {
				p.resolve(x)
			}
			x = p.parseIndexOrSlice(p.checkExpr(x))
		case token.LPAREN:
			if lhs {
				p.resolve(x)
			}
			x = p.parseCallOrConversion(p.checkExprOrType(x))
		case token.LBRACE:
			if isLiteralType(x) && (p.exprLev >= 0 || !isTypeName(x)) {
				if lhs {
					p.resolve(x)
				}
				x = p.parseLiteralValue(x)
			} else {
				break L
			}
		default:
			break L
		}
		lhs = false // no need to try to resolve again
	}

	return x
}

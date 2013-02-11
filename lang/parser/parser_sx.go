package parser

import (
	"simplex.sh/lang/ast"
	"simplex.sh/lang/token"
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

func (p *parser) parseTableType() *ast.TableType {
	if p.trace {
		defer un(trace(p, "TableType"))
	}

	pos := p.expect(token.TABLE)
	p.expect(token.LBRACK)
	key := p.parseType()
	p.expect(token.RBRACK)
	value := p.parseType()

	return &ast.TableType{Table: pos, Key: key, Value: value}
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
	case token.CHAN, token.ARROW:
		return p.parseChanType()
	case token.LPAREN:
		lparen := p.pos
		p.next()
		typ := p.parseType()
		rparen := p.expect(token.RPAREN)
		return &ast.ParenExpr{Lparen: lparen, X: typ, Rparen: rparen}
	}

	//=== start custom
	if p.mode&SimplexExtentions > 0 {
		switch p.tok {
		case token.VIEW:
			return p.parseViewType()
		case token.TABLE:
			return p.parseTableType()
		}
	}
	//=== end custom

	// no type found
	return nil
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

		//=== start custom
	case *ast.ViewType:
	case *ast.TableType:
		//=== end custom

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

				//=== start custom
			case token.SELECT:
				if p.mode&SimplexExtentions > 0 {
					pos := p.pos
					sel := &ast.Ident{NamePos: pos, Name: "select"}
					p.next()
					x = &ast.SelectorExpr{X: p.checkExpr(x), Sel: sel}
				} else {
					break L
				}
				//=== end custom

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

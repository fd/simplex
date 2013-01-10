package ast

import (
	"fmt"
)

type Visitor interface {
	Visit(Node) Visitor
}

type Visitable interface {
	Visit(Visitor)
}

func Walk(v Visitor, n Node) {
	if n == nil {
		return
	}

	v = v.Visit(n)
	if v == nil {
		return
	}

	switch val := n.(type) {

	case *Template:
		for _, stmt := range val.Statements {
			Walk(v, stmt)
		}

	case *Block:
		Walk(v, val.Macro)
		Walk(v, val.Expression)
		Walk(v, val.Template)
		if val.ElseTemplate != nil {
			Walk(v, val.ElseTemplate)
		}

	case *Interpolation:
		Walk(v, val.Expression)

	case *Literal:
		return

	case *Comment:
		return

	case *Identifier:
		return

	case *Expression:
		return

	case Visitable:
		val.Visit(v)

	default:
		panic(fmt.Sprintf("unhandled ast type (%T)", n))

	}
}

type Skip struct {
	V Visitor
	N int
}

func (visitor *Skip) Visit(n Node) Visitor {
	visitor.N -= 1

	if visitor.N > 0 {
		return visitor
	}

	return visitor.V
}

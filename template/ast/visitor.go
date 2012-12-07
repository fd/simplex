package ast

type Visitor interface {
	Visit(Node) Visitor
}

func Walk(v Visitor, n Node) {
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
		Walk(v, val.Expression)
		Walk(v, val.Template)
		if val.ElseTemplate != nil {
			Walk(v, val.ElseTemplate)
		}

	case *Interpolation:
		Walk(v, val.Expression)

	case *Comment:
		return

	case *Literal:
		return

	case *Get:
		Walk(v, val.From)
		Walk(v, val.Name)

	case *FunctionCall:
		Walk(v, val.From)

		for _, arg := range val.Args {
			Walk(v, arg)
		}

		for _, arg := range val.Options {
			Walk(v, arg)
		}

	case *IntegerLiteral:
		return

	case *FloatLiteral:
		return

	case *StringLiteral:
		return

	case *Identifier:
		return

	default:
		panic("unhandled type")

	}
}

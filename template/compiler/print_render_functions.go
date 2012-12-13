package compiler

import (
	"bytes"
	"fmt"
	"github.com/fd/w/template/ast"
	"strconv"
	"strings"
)

func (ctx *Context) PrintRenderFunctions() {
	for _, render := range ctx.RenderFuncs {
		var buf bytes.Buffer

		ast.Walk(&print_render_function{

			name: render.FunctionName(),
			buf:  &buf,
		}, render.Template)

		render.Golang = buf.String()
	}
}

type print_render_function struct {
	name string
	buf  *bytes.Buffer

	value_stack []string
	value_id    int
}

func (visitor *print_render_function) Visit(n ast.Node) ast.Visitor {

	switch v := n.(type) {
	case *ast.Template:
		visitor.print_debug_info(n)
		visitor.printf(
			"func %s(ctx Context, val Value)Value{\nvar buf bytes.Buffer\n",
			visitor.name,
		)

		ast.Walk(&ast.Skip{visitor, 1}, n)

		visitor.print_debug_info(n)
		visitor.printf(
			"return buf.String()\n}",
		)

		return nil

	case *ast.Interpolation:
		ast.Walk(&ast.Skip{visitor, 1}, n)

		//val := visitor.pop_value_stack()

		visitor.print_debug_info(n)
		visitor.printf(
			"runtime.WriteHtml(&buf, %s)\n",
			//      val,
			"{{stack}}",
		)
		return nil

	case *ast.Comment:
		visitor.print_debug_info(n)
		visitor.printf(
			"runtime.WriteHtml(&buf, runtime.HTML(%s))\n",
			strconv.Quote("<!-- "+v.Content+" -->"),
		)
		return nil

	case *ast.Literal:
		visitor.print_debug_info(n)
		visitor.printf(
			"runtime.WriteHtml(&buf, runtime.HTML(%s))\n",
			strconv.Quote(v.Content),
		)
		return nil

	case *ast.IntegerLiteral:
		val := visitor.push_value_stack()

		visitor.print_debug_info(n)
		visitor.printf(
			"var %s int = %d\n",
			val,
			v.Value,
		)
		return nil

	case *ast.FloatLiteral:
		val := visitor.push_value_stack()

		visitor.print_debug_info(n)
		visitor.printf(
			"var %s float64 = %f\n",
			val,
			v.Value,
		)
		return nil

	case *ast.StringLiteral:
		val := visitor.push_value_stack()

		visitor.print_debug_info(n)
		visitor.printf(
			"var %s string = %s\n",
			val,
			strconv.Quote(v.Value),
		)
		return nil

	default:
		visitor.print_debug_info(n)
		visitor.printf("/*\n  Unhandled:\n    %s\n*/\n", strings.Replace(n.String(), "\n", "\n    ", -1))

	}

	return visitor
}

func (visitor *print_render_function) print_debug_info(n ast.Node) {
	i := n.NodeInfo()
	visitor.printf(
		"/*\n  file:   %s\n  line:   %d\n  column: %d\n*/\n",
		i.File, i.Line, i.Column,
	)
}

func (visitor *print_render_function) printf(f string, a ...interface{}) {
	fmt.Fprintf(visitor.buf, f, a...)
}

func (visitor *print_render_function) push_value_stack() string {
	visitor.value_id += 1
	value_name := fmt.Sprintf("value_%d", visitor.value_id)
	visitor.value_stack = append(visitor.value_stack, value_name)
	return value_name
}

func (visitor *print_render_function) pop_value_stack() string {
	if l := len(visitor.value_stack); l > 0 {
		val := visitor.value_stack[l-1]
		visitor.value_stack = visitor.value_stack[:l-1]
		return val
	}
	panic(fmt.Sprintf("Unable to pop empty value stack."))
}

package compiler

import (
	"bytes"
	"fmt"
	"github.com/fd/w/template/ast"
	"sort"
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
	defer func() {
		if e := recover(); e != nil {
			fmt.Printf("e: %s\n", e)
			fmt.Printf("n: %s\n", n)
			panic(e)
		}
	}()

	switch v := n.(type) {
	case *ast.Template:
		visitor.printf(
			"func %s(ctx Context, val Value)Value{\n",
			visitor.name,
		)
		visitor.print_debug_info(n)
		visitor.printf(
			"var buf bytes.Buffer\n",
		)

		ast.Walk(&ast.Skip{visitor, 1}, n)

		visitor.print_debug_info(n)
		visitor.printf(
			"return buf.String()\n}\n\n",
		)

		return nil

	case *ast.Interpolation:
		ast.Walk(&ast.Skip{visitor, 1}, n)

		val := visitor.pop_value_stack()

		visitor.print_debug_info(n)
		visitor.printf(
			"runtime.WriteHtml(&buf, %s)\n",
			val,
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

	case *ast.Identifier:
		if v.Value == "$" {
			val := visitor.push_value_stack()

			visitor.print_debug_info(n)
			visitor.printf(
				"var %s = val\n",
				val,
			)
		}

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

	case *ast.Get:
		ast.Walk(&ast.Skip{visitor, 1}, n)

		val1 := visitor.pop_value_stack()
		val2 := visitor.push_value_stack()

		visitor.print_debug_info(n)
		visitor.printf(
			"%s := Get(%s, %s)\n",
			val2,
			val1,
			strconv.Quote(v.Name.Value),
		)
		return nil

	case *ast.FunctionCall:
		args := []string{}

		if v.From != nil {
			ast.Walk(visitor, v.From)
			val := visitor.pop_value_stack()
			args = append(args, val)
		}

		for _, a := range v.Args {
			ast.Walk(visitor, a)
			val := visitor.pop_value_stack()
			args = append(args, val)
		}

		if len(v.Options) > 0 {
			opts_pairs := []string{}

			for n, a := range v.Options {
				ast.Walk(visitor, a)
				val := visitor.pop_value_stack()
				opts_pairs = append(opts_pairs, fmt.Sprintf("%s: %s", strconv.Quote(n), val))
			}

			sort.Strings(opts_pairs)

			res := visitor.push_value_stack()
			opts_str := strings.Join(opts_pairs, ", ")

			visitor.print_debug_info(n)
			visitor.printf(
				"%s := map[string]interface{}{ %s }\n",
				res,
				opts_str,
			)
			args = append(args, res)
		}

		res := visitor.push_value_stack()
		args_str := strings.Join(args, ", ")
		visitor.print_debug_info(n)
		visitor.printf(
			"%s := %s(%s)\n",
			res,
			v.Name,
			args_str,
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

	fmt.Printf(visitor.buf.String())

	panic(fmt.Sprintf("Unable to pop empty value stack."))
}

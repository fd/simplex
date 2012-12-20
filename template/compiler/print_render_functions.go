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
			render_func: render,
			ctx:         ctx,
			name:        render.FunctionName(),
			buf:         &buf,
		}, render.Template)

		render.Golang = buf.String()
	}
}

type print_render_function struct {
	ctx         *Context
	render_func *RenderFunc

	name string
	buf  *bytes.Buffer

	value_stack []string
	value_id    int
}

func (visitor *print_render_function) Visit(n ast.Node) ast.Visitor {
	imports := visitor.ctx.ImportsFor(visitor.render_func.ImportPath)

	switch v := n.(type) {
	case *ast.Template:
		data_pkg := imports.Register("github.com/fd/w/data")
		runtime_pkg := imports.Register("github.com/fd/w/runtime")

		visitor.printf(
			"func %s(ctx %s.Context, val %s.Value) *%s.Buffer {\n",
			visitor.name,
			data_pkg, data_pkg, runtime_pkg,
		)
		visitor.print_debug_info(n)
		visitor.printf(
			"buf := new(%s.Buffer)\n",
			runtime_pkg,
		)

		ast.Walk(&ast.Skip{visitor, 1}, n)

		visitor.print_debug_info(n)
		visitor.printf(
			"return buf\n}\n\n",
		)

		return nil

	case *ast.Interpolation:
		ast.Walk(&ast.Skip{visitor, 1}, n)

		val := visitor.pop_value_stack()

		visitor.print_debug_info(n)
		visitor.printf(
			"buf.Write(%s)\n",
			val,
		)
		return nil

	case *ast.Comment:
		runtime_pkg := imports.Register("github.com/fd/w/runtime")

		visitor.print_debug_info(n)
		visitor.printf(
			"buf.Write(%s.HTML(%s))\n",
			runtime_pkg,
			strconv.Quote("<!-- "+v.Content+" -->"),
		)
		return nil

	case *ast.Literal:
		visitor.print_debug_info(n)
		s := v.Content
		for len(s) > 0 {
			runtime_pkg := imports.Register("github.com/fd/w/runtime")
			c := ""
			if len(s) >= 100 {
				c = s[:100]
				s = s[100:]
			} else {
				c = s[:len(s)]
				s = ""
			}
			visitor.printf(
				"buf.Write(%s.HTML(%s))\n",
				runtime_pkg,
				strconv.Quote(c),
			)
		}
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

		value_pkg := imports.Register("github.com/fd/w/data/value")
		visitor.printf(
			"var %s = %s.Get(%s, %s)\n",
			val2,
			value_pkg,
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

package compiler

import (
	"fmt"
	"github.com/fd/w/template/ast"
	"path"
	"unicode"
)

func (ctx *Context) LookupFunctionCalls() {
	for _, render := range ctx.RenderFuncs {
		ast.Walk(&lookup_function_calls{ctx, render.ImportPath}, render.Template)
	}
}

type lookup_function_calls struct {
	ctx         *Context
	import_path string
}

func (visitor *lookup_function_calls) Visit(n ast.Node) ast.Visitor {

	switch v := n.(type) {
	case *ast.FunctionCall:
		ast.Walk(&ast.Skip{visitor, 1}, v)
		visitor.Perform(v)
		return nil
	}

	return visitor
}

func (visitor *lookup_function_calls) Perform(fc *ast.FunctionCall) {

	// CASE 1: util.Escape()
	// when fc.From is a lowercase Literal look for a function matching the profile
	//   in $GOPATH/src/example/apps/example_app/module/
	//   util.Escape()
	//   > "example/apps/example_app/module/util".Escape()
	//   > "example/apps/example_app/util".Escape()

	// CASE 2: Escape()
	// when fc.From is nil look for a function matching the profile
	//   in $GOPATH/src/example/apps/example_app/module/
	//   Escape()
	//   > helpers.Escape()
	//   > "example/apps/example_app/module/helpers".Escape()
	//   > "example/apps/example_app/helpers".Escape()

	// CASE 3: escape()
	// when fc.From is nil look for a function matching the profile
	//   in $GOPATH/src/example/apps/example_app/module/
	//   escape()
	//   > module.escape()
	//   > "example/apps/example_app/module".escape()

	// CASE 1:
	if l, is_lit := fc.From.(*ast.Identifier); is_lit && l != nil && is_lower(l.Value) && is_upper(fc.Name) {
		name := visitor.resolve_function_name(l.Value, fc.Name)
		if name != "" {
			fc.From = nil
			fc.Name = name
			return
		}
	}

	// CASE 2:
	if fc.From == nil && is_upper(fc.Name) {
		name := visitor.resolve_function_name("helpers", fc.Name)
		if name != "" {
			fc.From = nil
			fc.Name = name
			return
		}
	}

	// CASE 3:
	if fc.From == nil && is_lower(fc.Name) {
		name := visitor.lookup_function_name(".", fc.Name)
		if name != "" {
			fc.From = nil
			fc.Name = name
			return
		}
	}

}

func (visitor *lookup_function_calls) resolve_function_name(pkg_name, func_name string) string {
	base := visitor.import_path

	for {
		pkg_path := path.Join(base, pkg_name)
		fullname := fmt.Sprintf("\"%s\".%s", pkg_path, func_name)

		if _, p := visitor.ctx.Helpers[fullname]; p {
			i := visitor.ctx.ImportsFor(visitor.import_path)
			name := i.Register(pkg_path)
			return name + "." + func_name
		}

		base = path.Dir(base)

		if base == "." || base == "/" {
			break
		}
	}

	return ""
}

func (visitor *lookup_function_calls) lookup_function_name(pkg_name, func_name string) string {
	base := visitor.import_path
	pkg_path := path.Join(base, pkg_name)
	fullname := fmt.Sprintf("\"%s\".%s", pkg_path, func_name)

	if _, p := visitor.ctx.Helpers[fullname]; p {
		i := visitor.ctx.ImportsFor(visitor.import_path)
		name := i.Register(pkg_path)
		return name + "." + func_name
	}

	return ""
}

func is_lower(s string) bool {
	if len(s) == 0 {
		return false
	}
	return unicode.IsLower(rune(s[0])) && unicode.IsLetter(rune(s[0]))
}

func is_upper(s string) bool {
	if len(s) == 0 {
		return false
	}
	return unicode.IsUpper(rune(s[0])) && unicode.IsLetter(rune(s[0]))
}

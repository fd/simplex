package compiler

import (
	"fmt"
	"github.com/fd/simplex/ast"
	"github.com/fd/simplex/printer"
	"github.com/fd/simplex/token"
	"github.com/fd/simplex/types"
	"io"
	"os"
	"sort"
	"text/template"
)

type printer_t struct {
	ctx                 *Context
	reflect_import_name string
	runtime_import_name string
}

func (c *Context) print_go() error {
	var w io.Writer

	f, err := os.Create(c.OutputFile)
	if err != nil {
		return err
	}
	defer f.Close()

	w = f
	w = io.MultiWriter(w, os.Stdout)

	merged_sx_file := c.merge_simplex_files()

	p := &printer_t{ctx: c}

	err = p.print_intro_and_merged_sx_file(w, c.TypesPackage.Name, merged_sx_file)
	if err != nil {
		return err
	}

	// print table interfaces
	// print table structs
	err = p.print_tables(w, c.TableTypes)
	if err != nil {
		return err
	}

	// print keyed view interfaces
	// print keyed view structs
	err = p.print_keyed_views(w, c.ViewTypes)
	if err != nil {
		return err
	}

	// print indexed view interfaces
	// print indexed view structs
	err = p.print_indexed_views(w, c.ViewTypes)
	if err != nil {
		return err
	}

	return nil
}

func (c *Context) merge_simplex_files() *ast.File {
	sx_files := make(map[string]*ast.File, len(c.SxFiles))
	for _, name := range c.SxFiles {
		sx_files[name] = c.AstFiles[name]
	}

	pkg := &ast.Package{
		Name:  c.TypesPackage.Name,
		Files: sx_files,
	}
	file := ast.MergePackageFiles(pkg, ast.FilterImportDuplicates)
	collect_imports_at_the_top(file)

	return file
}

func collect_imports_at_the_top(f *ast.File) {
	decl := f.Decls
	imports := []ast.Decl{}

	for i := len(decl) - 1; i >= 0; i-- {
		n := decl[i]
		d, ok := n.(*ast.GenDecl)
		if !ok || d.Tok != token.IMPORT {
			continue
		}

		imports = append(imports, d)

		if i > 0 {
			decl = append(decl[:i], decl[i+1:]...)
		} else {
			decl = decl[i+1:]
		}
	}

	f.Decls = append(imports, decl...)
}

func (p *printer_t) print_intro_and_merged_sx_file(w io.Writer, pkg_name string, sx_file *ast.File) error {
	var (
		reflect_import_name string
		runtime_import_name string
	)

	for _, spec := range sx_file.Imports {
		if spec.Path.Kind == token.STRING {
			switch spec.Path.Value {

			case `"reflect"`:
				reflect_import_name = "reflect"
				if spec.Name != nil {
					reflect_import_name = spec.Name.Name
				}

			case `"github.com/fd/simplex/runtime"`:
				runtime_import_name = "runtime"
				if spec.Name != nil {
					runtime_import_name = spec.Name.Name
				}

			}
		}
	}

	if reflect_import_name == "" {
		reflect_import_name = "sx_reflect"

		imp := &ast.ImportSpec{
			Name: ast.NewIdent("sx_reflect"),
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: `"reflect"`,
			},
		}

		decl := &ast.GenDecl{
			Tok:   token.IMPORT,
			Specs: []ast.Spec{imp},
		}

		sx_file.Imports = append(sx_file.Imports, imp)
		sx_file.Decls = append([]ast.Decl{decl}, sx_file.Decls...)
	}

	if runtime_import_name == "" {
		runtime_import_name = "sx_runtime"

		imp := &ast.ImportSpec{
			Name: ast.NewIdent("sx_runtime"),
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: `"github.com/fd/simplex/runtime"`,
			},
		}

		decl := &ast.GenDecl{
			Tok:   token.IMPORT,
			Specs: []ast.Spec{imp},
		}

		sx_file.Imports = append(sx_file.Imports, imp)
		sx_file.Decls = append([]ast.Decl{decl}, sx_file.Decls...)
	}

	p.reflect_import_name = reflect_import_name
	p.runtime_import_name = runtime_import_name

	conf := printer.Config{Mode: printer.SourcePos, Tabwidth: 8}
	err := conf.Fprint(w, p.ctx.FileSet, sx_file)
	if err != nil {
		return err
	}

	fmt.Fprintln(w, "\n//line sx_generated.go:1")

	return nil
}

var table_tmpl = template.Must(template.New("table_tmpl").Parse(`
type (
  {{.TypeName}} interface {
    {{.Runtime}}.GenericTable
    {{.ViewTypeName}}
  }

  sx_{{.TypeName}} struct {
    sx_{{.ViewTypeName}}
  }
)
func (s sx_{{.TypeName}}) InnerTable() *{{.Runtime}}.Table { return s.Table }
func new_{{.TypeName}}() {{.TypeName}} { return sx_{{.TypeName}}{} }
func wrap_{{.TypeName}}(tab *{{.Runtime}}.Table) {{.TypeName}} {
  t := sx_{{.TypeName}}{}
  t.Table = tab
  return t
}
`))

func (p *printer_t) print_tables(w io.Writer, tables map[string]*types.Table) error {
	type data struct {
		Runtime string
		Reflect string

		TypeName string
		KeyType  string
		EltType  string
		KeyZero  string
		EltZero  string

		ViewTypeName string
	}

	names := make([]string, 0, len(tables))
	for name := range tables {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		typ := tables[name]
		err := table_tmpl.Execute(w, data{
			Runtime: p.runtime_import_name,
			Reflect: p.reflect_import_name,

			TypeName: name,
			KeyType:  type_string(typ.Key),
			EltType:  type_string(typ.Elt),
			KeyZero:  type_zero(typ.Key),
			EltZero:  type_zero(typ.Elt),

			ViewTypeName: type_string(&types.View{typ.Key, typ.Elt}),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

var keyed_view_tmpl = template.Must(template.New("keyed_view_tmpl").Parse(`
type (
  {{.TypeName}} interface {
    {{.Runtime}}.GenericKeyedView
    {{.IndexedTypeName}}
    KeyZero() {{.KeyType}}
  }

  sx_{{.TypeName}} struct {
    sx_{{.IndexedTypeName}}
  }
)
func (s sx_{{.TypeName}}) KeyType() {{.Reflect}}.Type { return {{.Reflect}}.TypeOf(s.KeyZero()) }
func (s sx_{{.TypeName}}) KeyZero() {{.KeyType}} { return {{.KeyZero}} }
func wrap_{{.TypeName}}(tab *{{.Runtime}}.Table) {{.TypeName}} {
  t := sx_{{.TypeName}}{}
  t.Table = tab
  return t
}
`))

func (p *printer_t) print_keyed_views(w io.Writer, views map[string]*types.View) error {
	type data struct {
		Runtime string
		Reflect string

		TypeName string
		KeyType  string
		EltType  string
		KeyZero  string
		EltZero  string

		IndexedTypeName string
	}

	names := make([]string, 0, len(views))
	for name := range views {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		typ := views[name]
		if typ.Key == nil {
			continue
		}

		err := keyed_view_tmpl.Execute(w, &data{
			Runtime: p.runtime_import_name,
			Reflect: p.reflect_import_name,

			TypeName: name,
			KeyType:  type_string(typ.Key),
			EltType:  type_string(typ.Elt),
			KeyZero:  type_zero(typ.Key),
			EltZero:  type_zero(typ.Elt),

			IndexedTypeName: type_string(&types.View{nil, typ.Elt}),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

var indexed_view_tmpl = template.Must(template.New("indexed_view_tmpl").Parse(`
type (
  {{.TypeName}} interface {
    {{.Runtime}}.GenericView
    {{.Runtime}}.GenericIndexedView
    EltZero() {{.EltType}}
  }

  sx_{{.TypeName}} struct {
    *{{.Runtime}}.Table
  }
)
func wrap_{{.TypeName}}(tab *{{.Runtime}}.Table) {{.TypeName}} { return sx_{{.TypeName}}{ Table: tab } }
func (s sx_{{.TypeName}}) EltType() {{.Reflect}}.Type { return {{.Reflect}}.TypeOf(s.EltZero()) }
func (s sx_{{.TypeName}}) EltZero() {{.EltType}} { return {{.EltZero}} }
func (s sx_{{.TypeName}}) InnerView() *{{.Runtime}}.Table { return s.Table }
`))

func (p *printer_t) print_indexed_views(w io.Writer, views map[string]*types.View) error {
	type data struct {
		Runtime string
		Reflect string

		TypeName string
		EltType  string
		EltZero  string
	}

	names := make([]string, 0, len(views))
	for name := range views {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		typ := views[name]
		if typ.Key != nil {
			continue
		}

		err := indexed_view_tmpl.Execute(w, data{
			Runtime: p.runtime_import_name,
			Reflect: p.reflect_import_name,

			TypeName: name,
			EltType:  type_string(typ.Elt),
			EltZero:  type_zero(typ.Elt),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

package compiler

import (
	"github.com/fd/simplex/printer"
	"github.com/fd/simplex/types"
	"io"
	"os"
	"path"
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

	conf := printer.Config{Mode: printer.SourcePos, Tabwidth: 8}
	for _, name := range c.SxFiles {
		f, err := os.Create(path.Join(c.OutputDir, name[:len(name)-3]+".go"))
		if err != nil {
			return err
		}
		defer f.Close()

		w = f
		w = io.MultiWriter(w, os.Stdout)

		err = conf.Fprint(w, c.FileSet, c.AstFiles[name])
		if err != nil {
			return err
		}
	}

	f, err := os.Create(path.Join(c.OutputDir, "smplx_generated.go"))
	if err != nil {
		return err
	}
	defer f.Close()

	w = f
	w = io.MultiWriter(w, os.Stdout)

	p := &printer_t{ctx: c}

	err = p.print_intro(w, c.TypesPackage.Name)
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

var intro_tmpl = template.Must(template.New("intro_tmpl").Parse(`
package {{.PkgName}}

import (
  sx_reflect "reflect"
  sx_runtime "github.com/fd/simplex/runtime"
)

`))

func (p *printer_t) print_intro(w io.Writer, pkg_name string) error {
	p.reflect_import_name = "sx_reflect"
	p.runtime_import_name = "sx_runtime"

	type data struct {
		PkgName string
	}

	return intro_tmpl.Execute(w, data{pkg_name})
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

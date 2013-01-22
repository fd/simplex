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
	printed_types       map[string]bool
}

func (c *Context) print_go() error {
	err := c.print_sx_files_as_go_files()
	if err != nil {
		return err
	}

	err = c.print_generated_go()
	if err != nil {
		return err
	}

	return nil
}

func (c *Context) print_generated_go() error {
	var w io.Writer

	f, err := os.Create(path.Join(c.OutputDir, "smplx_generated.go"))
	if err != nil {
		return err
	}
	defer f.Close()

	w = f
	w = io.MultiWriter(w, os.Stdout)

	p := &printer_t{ctx: c, printed_types: map[string]bool{}}

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

func (c *Context) print_sx_files_as_go_files() error {
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
    {{.ViewTypeName}}
    TableId() string
  }

  sx_{{.TypeName}} struct {
    sx_{{.ViewTypeName}}
    Id string
  }
)
func (t sx_{{.TypeName}}) TableId() string { return t.Id }
func (t sx_{{.TypeName}}) Resolve(txn *{{.Runtime}}.Transaction, events chan<- {{.Runtime}}.Event) {
  t.Deferred.Resolve(txn, events)
}
func new_{{.TypeName}}(env *{{.Runtime}}.Environment, id string) {{.TypeName}} {
  t := sx_{{.TypeName}}{}
  t.Id = id
  t.Deferred = {{.Runtime}}.DeclareTable(t.Id)
  env.RegisterTable(t)
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
		if p.printed_types[name] {
			continue
		}
		p.printed_types[name] = true

		typ := tables[name]
		parent := &types.View{typ.Key, typ.Elt}
		p.ctx.ViewTypes[go_type_string(parent)] = parent

		err := table_tmpl.Execute(w, data{
			Runtime: p.runtime_import_name,
			Reflect: p.reflect_import_name,

			TypeName: name,
			KeyType:  go_type_string(typ.Key),
			EltType:  go_type_string(typ.Elt),
			KeyZero:  type_zero(typ.Key),
			EltZero:  type_zero(typ.Elt),

			ViewTypeName: go_type_string(parent),
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
    {{.IndexedTypeName}}
    KeyType() sx_reflect.Type
    KeyZero() {{.KeyType}}
  }

  sx_{{.TypeName}} struct {
    sx_{{.IndexedTypeName}}
  }
)
func (s sx_{{.TypeName}}) KeyType() {{.Reflect}}.Type { return {{.Reflect}}.TypeOf(s.KeyZero()) }
func (s sx_{{.TypeName}}) KeyZero() {{.KeyType}} { return {{.KeyZero}} }
func (t sx_{{.TypeName}}) Resolve(txn *{{.Runtime}}.Transaction, events chan<- {{.Runtime}}.Event) {
  t.Deferred.Resolve(txn, events)
}
func wrap_{{.TypeName}}(def {{.Runtime}}.Deferred) {{.TypeName}} {
  t := sx_{{.TypeName}}{}
  t.Deferred = def
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

		if p.printed_types[name] {
			continue
		}
		p.printed_types[name] = true

		parent := &types.View{nil, typ.Elt}
		p.ctx.ViewTypes[go_type_string(parent)] = parent

		err := keyed_view_tmpl.Execute(w, &data{
			Runtime: p.runtime_import_name,
			Reflect: p.reflect_import_name,

			TypeName: name,
			KeyType:  go_type_string(typ.Key),
			EltType:  go_type_string(typ.Elt),
			KeyZero:  type_zero(typ.Key),
			EltZero:  type_zero(typ.Elt),

			IndexedTypeName: go_type_string(parent),
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
    EltType() sx_reflect.Type
    EltZero() {{.EltType}}
    Resolve(txn *{{.Runtime}}.Transaction, events chan<- {{.Runtime}}.Event)
  }

  sx_{{.TypeName}} struct {
    Deferred {{.Runtime}}.Deferred
  }
)
func (s sx_{{.TypeName}}) EltType() {{.Reflect}}.Type { return {{.Reflect}}.TypeOf(s.EltZero()) }
func (s sx_{{.TypeName}}) EltZero() {{.EltType}} { return {{.EltZero}} }
func (t sx_{{.TypeName}}) Resolve(txn *{{.Runtime}}.Transaction, events chan<- {{.Runtime}}.Event) {
  t.Deferred.Resolve(txn, events)
}
func wrap_{{.TypeName}}(def {{.Runtime}}.Deferred) {{.TypeName}} {
  t := sx_{{.TypeName}}{}
  t.Deferred = def
  return t
}
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

		if p.printed_types[name] {
			continue
		}
		p.printed_types[name] = true

		err := indexed_view_tmpl.Execute(w, data{
			Runtime: p.runtime_import_name,
			Reflect: p.reflect_import_name,

			TypeName: name,
			EltType:  go_type_string(typ.Elt),
			EltZero:  type_zero(typ.Elt),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

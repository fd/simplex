package compiler

import (
	"github.com/fd/simplex/ast"
	"github.com/fd/simplex/token"
	"github.com/fd/simplex/types"
	"io"
	"os"
	"sort"
	"text/template"
)

type Context struct {
	OutputFile string
	GoFiles    []string
	SxFiles    []string

	AstFiles     map[string]*ast.File
	TypesPackage *types.Package
	ViewTypes    map[string]*types.View
	TableTypes   map[string]*types.Table
	FileSet      *token.FileSet
}

func (c *Context) Compile() error {
	var err error

	err = c.parse_files()
	if err != nil {
		return err
	}

	err = c.check_types()
	if err != nil {
		return err
	}

	err = c.generate_go()
	if err != nil {
		return err
	}

	return nil
}

func (c *Context) generate_go() error {
	var w io.Writer

	f, err := os.Create(c.OutputFile)
	if err != nil {
		return err
	}
	defer f.Close()

	w = f
	w = io.MultiWriter(w, os.Stdout)

	err = print_intro(w, c.TypesPackage.Name)
	if err != nil {
		return err
	}

	// print table interfaces
	// print table structs
	err = print_tables(w, c.TableTypes)
	if err != nil {
		return err
	}

	// print keyed view interfaces
	// print keyed view structs
	err = print_keyed_views(w, c.ViewTypes)
	if err != nil {
		return err
	}

	// print indexed view interfaces
	// print indexed view structs
	err = print_indexed_views(w, c.ViewTypes)
	if err != nil {
		return err
	}

	// parse generated go file
	// replace type expr in .sx files
	// replace view methods in .sx files
	// merge .sx files into generated go file
	// print generated go file

	return nil
}

var intro_tmpl = template.Must(template.New("").Parse(`
package {{.PkgName}}

import (
  sx_runtime "github.com/fd/simplex/runtime"
  "reflect"
)
`))

func print_intro(w io.Writer, pkg_name string) error {
	type data struct {
		PkgName string
	}
	return intro_tmpl.Execute(w, data{
		PkgName: pkg_name,
	})
}

var table_tmpl = template.Must(template.New("table_tmpl").Parse(`
type (
  {{.TypeName}} interface {
    sx_runtime.GenericTable
    KeyZero() {{.KeyType}}
    EltZero() {{.EltType}}
  }

  sx_{{.TypeName}} struct {
    *sx_runtime.Table
  }
)
func (s *sx_{{.TypeName}}) KeyType() reflect.Type { return reflect.TypeOf(s.KeyZero()) }
func (s *sx_{{.TypeName}}) EltType() reflect.Type { return reflect.TypeOf(s.EltZero()) }
func (s *sx_{{.TypeName}}) KeyZero() {{.KeyType}} { return {{.KeyZero}} }
func (s *sx_{{.TypeName}}) EltZero() {{.EltType}} { return {{.EltZero}} }
func (s *sx_{{.TypeName}}) InnerTable() *sx_runtime.Table { return s.Table }
func (s *sx_{{.TypeName}}) InnerView() *sx_runtime.Table { return s.Table }
`))

func print_tables(w io.Writer, tables map[string]*types.Table) error {
	type data struct {
		TypeName string
		KeyType  string
		EltType  string
		KeyZero  string
		EltZero  string
	}

	names := make([]string, 0, len(tables))
	for name := range tables {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		typ := tables[name]
		err := table_tmpl.Execute(w, data{
			TypeName: name,
			KeyType:  type_string(typ.Key),
			EltType:  type_string(typ.Elt),
			KeyZero:  type_zero(typ.Key),
			EltZero:  type_zero(typ.Elt),
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
    sx_runtime.GenericKeyedView
    KeyZero() {{.KeyType}}
    EltZero() {{.EltType}}
  }

  sx_{{.TypeName}} struct {
    *sx_runtime.Table
  }
)
func (s *sx_{{.TypeName}}) KeyType() reflect.Type { return reflect.TypeOf(s.KeyZero()) }
func (s *sx_{{.TypeName}}) EltType() reflect.Type { return reflect.TypeOf(s.EltZero()) }
func (s *sx_{{.TypeName}}) KeyZero() {{.KeyType}} { return {{.KeyZero}} }
func (s *sx_{{.TypeName}}) EltZero() {{.EltType}} { return {{.EltZero}} }
func (s *sx_{{.TypeName}}) InnerView() *sx_runtime.Table { return s.Table }
`))

func print_keyed_views(w io.Writer, views map[string]*types.View) error {
	type data struct {
		TypeName string
		KeyType  string
		EltType  string
		KeyZero  string
		EltZero  string
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
			TypeName: name,
			KeyType:  type_string(typ.Key),
			EltType:  type_string(typ.Elt),
			KeyZero:  type_zero(typ.Key),
			EltZero:  type_zero(typ.Elt),
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
    sx_runtime.GenericIndexedView
    EltZero() {{.EltType}}
  }

  sx_{{.TypeName}} struct {
    *sx_runtime.Table
  }
)
func (s *sx_{{.TypeName}}) EltType() reflect.Type { return reflect.TypeOf(s.EltZero()) }
func (s *sx_{{.TypeName}}) EltZero() {{.EltType}} { return {{.EltZero}} }
func (s *sx_{{.TypeName}}) InnerView() *sx_runtime.Table { return s.Table }
`))

func print_indexed_views(w io.Writer, views map[string]*types.View) error {
	type data struct {
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

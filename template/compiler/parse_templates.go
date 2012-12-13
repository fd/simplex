package compiler

import (
	"fmt"
	w_parser "github.com/fd/w/template/parser"
	"os"
	"path"
	"strings"
)

func (ctx *Context) ParseTemplates() error {
	var errs Errors

	for path, pkg := range ctx.go_build_packages {
		dir := pkg.Dir
		err := ctx.parse_templates_in_dir(path, dir)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

func (ctx *Context) parse_templates_in_dir(import_path, dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}

	entries, err := d.Readdir(0)
	if err != nil {
		return err
	}

	for _, fi := range entries {
		base := fi.Name()

		if !strings.HasSuffix(base, ".go.html") {
			continue
		}

		tmpl, err := w_parser.ParseFile(path.Join(dir, base))
		if err != nil {
			return err
		}

		base = base[:len(base)-8]

		name := fmt.Sprintf("\"%s\".%s", import_path, base)
		ctx.RenderFuncs[name] = &RenderFunc{
			Name:       base,
			ImportPath: import_path,
			Template:   tmpl,
			Export:     true,
		}
	}

	return nil
}

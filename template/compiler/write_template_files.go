package compiler

import (
	"fmt"
	"io/ioutil"
	"path"
	"sort"
	"strconv"
	"strings"
)

func (ctx *Context) WriteTemplateFiles() error {
	files := map[string][]string{}

	{ // append imports
		for pkg_name, imports := range ctx.Imports {
			i := []string{}

			for path, name := range imports.imports {
				s := fmt.Sprintf("  %s %s", name, strconv.Quote(path))
				i = append(i, s)
			}

			sort.Strings(i)

			str := fmt.Sprintf("import(\n%s\n)\n", strings.Join(i, "\n"))

			c := files[pkg_name]
			c = append(c, str)
			files[pkg_name] = c
		}
	}

	{ // append render functions
		keys := make([]string, 0, len(ctx.RenderFuncs))
		for k := range ctx.RenderFuncs {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, key := range keys {
			render_func := ctx.RenderFuncs[key]
			c := files[render_func.ImportPath]
			c = append(c, render_func.Golang)
			files[render_func.ImportPath] = c
		}
	}

	{ // prepend `package` statements
		for n, c := range files {
			base := path.Base(n)
			stmt := fmt.Sprintf("package %s\n\n", base)
			c = append([]string{stmt}, c...)
			files[n] = c
		}
	}

	{ // write files
		for n, c := range files {
			content := strings.Join(c, "\n")
			build_pkg := ctx.go_build_packages[n]
			filename := path.Join(build_pkg.Dir, "zzz_w_compiled.go")
			err := ioutil.WriteFile(filename, []byte(content), 0644)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

package compiler

import (
	"os"
	"path/filepath"
)

func (ctx *Context) ImportPackages() error {
	// walk all dirs in ./apps/
	// -> import all
	p := filepath.Join(ctx.WROOT, "apps")

	err := filepath.Walk(p, func(path string, fi os.FileInfo, err error) error {
		if p == path {
			return nil
		}

		if err != nil {
			return err
		}

		if !fi.IsDir() {
			return nil
		}

		build_pkg, _, err := ctx.go_import_package(".", path)
		if err != nil {
			return err
		}

		if filepath.Base(filepath.Dir(build_pkg.ImportPath)) == "apps" {
			ctx.Applications[build_pkg.ImportPath] = true
		}

		return nil
	})
	if err != nil {
		return err
	}

	// walk all dirs in ./services/
	// -> import all

	return nil
}

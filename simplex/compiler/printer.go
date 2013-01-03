package compiler

import (
	"go/printer"
	"os"
	"path"
)

func (pkg *Package) Print() error {
	ast_f, ok := pkg.AstPackage.Files["smplx_generated.go"]
	if !ok {
		return nil
	}

	name := path.Join(pkg.BuildPackage.Dir, "smplx_generated.go")
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()

	c := printer.Config{printer.TabIndent | printer.SourcePos, 8}
	c.Fprint(f, pkg.FileSet, ast_f)

	f.Close()
	mod_time := pkg.ModTimes["smplx_generated.go"]
	os.Chtimes(name, mod_time, mod_time)

	return nil
}

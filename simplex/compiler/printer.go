package compiler

import (
	"go/printer"
	"os"
	"path"
)

func (pkg *Package) Print() ([]string, error) {
	fs, err := pkg.print_deps()
	if err != nil {
		return fs, err
	}

	ast_f, ok := pkg.AstPackage.Files["smplx_generated.go"]
	if !ok {
		return fs, nil
	}

	mod_time := pkg.ModTimes["smplx_generated.go"]
	stat, err := os.Stat(pkg.TargetPath)
	if len(fs) == 0 && err == nil && !mod_time.After(stat.ModTime()) {
		return fs, nil
	}
	err = nil

	name := path.Join(pkg.BuildPackage.Dir, "smplx_generated.go")
	fs = append(fs, name)
	f, err := os.Create(name)
	if err != nil {
		return fs, err
	}
	defer f.Close()

	c := printer.Config{printer.TabIndent | printer.SourcePos, 8}
	c.Fprint(f, pkg.FileSet, ast_f)

	f.Close()
	os.Chtimes(name, mod_time, mod_time)

	return fs, nil
}

func (pkg *Package) print_deps() (deps []string, err error) {

	for _, dep := range pkg.Imports {
		f, err := dep.Print()
		deps = append(deps, f...)

		if err != nil {
			break
		}
	}

	return deps, err
}

package compiler

import (
	"go/build"
	"os"
	"path"
	"strings"
)

func Import(dir string, srcDir string) (*Package, error) {

	{ // lookup in memory
		build_pkg, err := Context.Import(dir, srcDir, build.FindOnly)

		if _, ok := err.(*build.NoGoError); ok {
			err = nil
		}

		if err != nil {
			return nil, err
		}

		if pkg, p := Packages[build_pkg.ImportPath]; p {
			return pkg, nil
		}
	}

	build_pkg, err := Context.Import(dir, srcDir, 0)

	if _, ok := err.(*build.NoGoError); ok {
		err = nil
	}

	if err != nil {
		return nil, err
	}

	pkg := &Package{BuildPackage: build_pkg}

	if pkg.BuildPackage.IsCommand() {
		pkg.TargetPath = path.Join(pkg.BuildPackage.BinDir, path.Base(pkg.BuildPackage.Dir))
	} else {
		pkg.TargetPath = pkg.BuildPackage.PkgObj
	}

	err = pkg.find_simplex_files()
	if err != nil {
		return nil, err
	}

	Packages[build_pkg.ImportPath] = pkg
	return pkg, nil
}

func ImportResolved(dir string, srcDir string) (*Package, error) {
	pkg, err := Import(dir, srcDir)
	if err != nil {
		return nil, err
	}

	err = pkg.ParseFiles()
	if err != nil {
		return nil, err
	}

	err = pkg.ResolvePackage()
	if err != nil {
		return nil, err
	}

	err = pkg.MergeGeneratedFiles()
	if err != nil {
		return nil, err
	}

	return pkg, nil
}

func (pkg *Package) find_simplex_files() error {
	f, err := os.Open(pkg.BuildPackage.Dir)
	if err != nil {
		return err
	}
	defer f.Close()

	fis, err := f.Readdir(0)
	if err != nil {
		return err
	}

	for _, fi := range fis {
		name := fi.Name()

		if !strings.HasSuffix(name, ".smplx") {
			continue
		}

		if strings.HasPrefix(name, "_") {
			continue
		}

		pkg.SimplexFiles = append(pkg.SimplexFiles, name)
	}

	return nil
}

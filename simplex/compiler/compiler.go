package compiler

import (
	"os"
	"os/exec"
)

func (pkg *Package) Compile() error {
	pkg.Print()

	cmd := exec.Command("go", "get", "-v", ".")
	cmd.Dir = pkg.BuildPackage.Dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

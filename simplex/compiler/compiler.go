package compiler

import (
	"os"
	"os/exec"
)

func (pkg *Package) Compile() error {
	fs, err := pkg.Print()
	defer func() {
		for _, n := range fs {
			os.Remove(n)
		}
	}()
	if err != nil {
		return err
	}

	cmd := exec.Command("go", "get", "-v", ".")
	cmd.Dir = pkg.BuildPackage.Dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

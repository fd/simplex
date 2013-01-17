package compiler

import (
	"fmt"
	"os"
	"os/exec"
)

func (pkg *Package) Compile() error {
	fs, err := pkg.Print()
	fmt.Println(fs)
	//defer remove_all_files(fs)
	if err != nil {
		return err
	}

	cmd := exec.Command("go", "get", "-v", ".")
	cmd.Dir = pkg.BuildPackage.Dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func remove_all_files(fs []string) {
	for _, n := range fs {
		os.Remove(n)
	}
}

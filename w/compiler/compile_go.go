package compiler

import (
	"os"
	"os/exec"
	template "simplex.sh/w/template/compiler"
)

func compile_go_files(ctx *template.Context) error {
	cmd := exec.Command("go", "get", ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

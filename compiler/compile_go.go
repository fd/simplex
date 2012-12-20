package compiler

import (
	template "github.com/fd/w/template/compiler"
	"os"
	"os/exec"
)

func compile_go_files(ctx *template.Context) error {
	cmd := exec.Command("go", "get", ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

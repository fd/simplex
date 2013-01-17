package main

import (
	"github.com/fd/options"
	"os"
)

func Run() {
	var spec = options.MustParse(`
  w-dev -p PORT
  w toolchain
  --
  --
  --
  build  build  Build the container
  `)

	opts := spec.MustInterpret(os.Args, os.Environ())

	if len(opts.Args) == 0 {
		spec.PrintUsageAndExit()
	}

	switch opts.Command {
	case "build":
		Build(opts.Args)
	default:
		spec.PrintUsageAndExit()
	}
}

func main() {
	Run()
}

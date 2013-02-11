package main

import (
	"fmt"
	"github.com/fd/options"
	"os"
	"simplex.sh/w/compiler"
)

func Build(args []string) {
	var spec = options.MustParse(`
  w-dev build
  w toolchain
  --
  --
  --
  *
  `)

	opts := spec.MustInterpret(args, os.Environ())

	if len(opts.Args) != 0 {
		spec.PrintUsageAndExit()
	}

	var err error

	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = compiler.Compile(pwd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

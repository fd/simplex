package main

import (
	"fmt"
	"github.com/fd/options"
	"github.com/fd/w/template/compiler"
	"os"
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

	// compile schemas           (to go files)
	// compile templates         (to go files)
	err = compile_templates(pwd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// compile container wrapper (to go files)
	// compile additional tests  (to go files)
	// compile container         (to binary)
	// run tests
}

func compile_templates(pwd string) error {
	ctx := compiler.NewContext(pwd)

	err := ctx.Compile()
	if err != nil {
		return err
	}

	fmt.Printf("Helpers:\n")
	for n := range ctx.Helpers {
		fmt.Printf("- %s\n", n)
	}

	fmt.Printf("DataViews:\n")
	for n := range ctx.DataViews {
		fmt.Printf("- %s\n", n)
	}

	fmt.Printf("RenderFuncs:\n")
	for n, r := range ctx.RenderFuncs {
		fmt.Printf("- %s\n", n)
		fmt.Printf("%s\n", r.Golang)
	}

	return nil
}

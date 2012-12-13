package main

import (
	"github.com/fd/options"
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

	// compile schemas           (to go files)
	// compile templates         (to go files)
	// compile container wrapper (to go files)
	// compile additional tests  (to go files)
	// compile container         (to binary)
	// run tests
}

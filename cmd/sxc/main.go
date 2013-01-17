package main

import (
	"flag"
	"fmt"
	"github.com/fd/simplex/compiler"
	"os"
	"strings"
)

var (
	output = flag.String("o", "", "path of generated go file")
)

func usage() {
	fmt.Fprintf(os.Stderr,
		"usage: sxc -o output [input ...]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	sxfiles := make([]string, 0, len(flag.Args()))
	gofiles := make([]string, 0, len(flag.Args()))
	for _, name := range flag.Args() {
		if strings.HasSuffix(name, ".go") {
			gofiles = append(gofiles, name)
			continue
		}
		if strings.HasSuffix(name, ".sx") {
			sxfiles = append(sxfiles, name)
			continue
		}
	}

	ctx := compiler.Context{
		OutputFile: *output,
		GoFiles:    gofiles,
		SxFiles:    sxfiles,
	}

	err := ctx.Compile()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

package main

import (
	"fmt"
	"github.com/fd/w/simplex/compiler"
)

func main() {
	pkg, err := compiler.ImportResolved(
		"github.com/fd/w/simplex/example", ".")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = pkg.Compile()
	if err != nil {
		fmt.Println(err)
		return
	}
}

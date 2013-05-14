package main

import (
	"simplex.sh/container"
	_ "simplex.sh/example"
	_ "simplex.sh/store/file"
)

func main() { container.CLI() }

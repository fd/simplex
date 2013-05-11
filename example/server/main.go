package main

import (
	"github.com/fd/static/container"
	_ "github.com/fd/static/example"
	_ "github.com/fd/static/store/file"
)

func main() { container.CLI() }

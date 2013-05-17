package main

import (
	"simplex.sh/container"
	_ "simplex.sh/example"
	_ "simplex.sh/store/file"
	_ "simplex.sh/store/redis"
)

func main() { container.CLI() }

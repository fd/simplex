package main

// Libraries
import (
	"github.com/fd/w/container"
)

// Apps
import (
	_ "github.com/fd/w/apps/cp/locations"
	_ "github.com/fd/w/apps/cp/partners"
)

func main() {
	container.Run()
}

package main

// Libraries
import (
	"github.com/fd/w/container"
)

// Apps
import (
	_ "github.com/fd/w/example/apps/orakel/models"
)

func main() {
	container.Run()
}

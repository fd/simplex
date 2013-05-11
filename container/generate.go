package container

import (
	"fmt"
	"runtime"
)

func Generate(env Environment) error {
	c, err := new_container(env)
	if err != nil {
		return err
	}

	return c.Generate()
}

// generate all the applications
func (c *container_t) Generate() error {
	fmt.Printf("Generating everything...\n")
	defer fmt.Printf("Done\n")

	runtime.GOMAXPROCS(runtime.NumCPU() * 2)

	for _, app := range c.apps {
		fmt.Printf("- %s...", app.Name)
		err := app.Generate()
		if err != nil {
			c.err.Add(err)
		}
		fmt.Printf(" done\n")
	}

	return c.err.Normalize()
}

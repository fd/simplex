package container

import (
	"fmt"
	"runtime"
	"time"
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
	start := time.Now()
	fmt.Printf("Generating everything...\n")
	defer func() {
		fmt.Printf("Done (duration=%s)\n", time.Now().Sub(start))
	}()

	runtime.GOMAXPROCS(runtime.NumCPU() * 2)

	for _, app := range c.apps {
		app_start := time.Now()
		fmt.Printf("- %s...", app.Name)

		err := app.Generate()
		if err != nil {
			c.err.Add(err)
		}

		fmt.Printf(" done (duration=%s)\n", time.Now().Sub(app_start))
	}

	return c.err.Normalize()
}

package container

import (
	"fmt"
	"net/http"
	"os/signal"
	"runtime"
	"syscall"
)

func Serve(env Environment) error {
	c, err := new_container(env)
	if err != nil {
		return err
	}

	return c.Serve()
}

// Run the container
func (c *container_t) Serve() error {
	fmt.Printf("Listening at %s...\n", c.env.HttpAddr)
	defer fmt.Printf("Shutting down\n")

	runtime.GOMAXPROCS(runtime.NumCPU() * 2)

	go c.go_serve()

	signal.Notify(c.shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-c.shutdown

	return c.err.Normalize()
}

func (c *container_t) go_serve() {
	defer func() { c.shutdown <- syscall.SIGTERM }()

	err := http.ListenAndServe(c.env.HttpAddr, c)
	c.err.Add(err)
}

func (c *container_t) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.router.ServeHTTP(w, r)
}

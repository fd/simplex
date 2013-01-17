package container

import (
	"fmt"
	"github.com/fd/options"
	"github.com/fd/simplex/w/data"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

import (
	_ "github.com/fd/simplex/w/data/storage/file_system"
)

var spec = options.MustParse(`
container -p PORT
Run the container
--
!port=    -p,PORT    The port to listen on
--
!source=  SOURCE     The source database
!state=   STATE      The state database
!target=  TARGET     The target database
--
*
`)

func Run() {
	opts := spec.MustInterpret(os.Args, os.Environ())

	if len(opts.Args) != 0 {
		spec.PrintUsageAndExit()
	}

	runtime.GOMAXPROCS(runtime.NumCPU() * 2)

	err := data.Setup(opts.Get("source"), opts.Get("state"), opts.Get("state"))
	if err != nil {
		spec.PrintUsageWithError(err)
	}

	data.Run()

	fmt.Printf("=======================================\n")
	fmt.Printf("== Hello there                       ==\n")

	c := data.Changes{}

	/*
	   c.Update = map[string]data.Value{
	     "hello": data.Object{
	       "_id":   "hello",
	       "_type": "location",
	       "name":  "Hello World!",
	     },
	     "hi": data.Object{
	       "_id":   "hi",
	       "_type": "location",
	       "name":  "Hi World!",
	     },
	   }
	*/

	///*
	c.Create = map[string]data.Value{
		"hello": data.Object{
			"_id":   "hello",
			"_type": "location",
			"name":  "Hello World",
		},
		"hi": data.Object{
			"_id":   "hi",
			"_type": "location",
			"name":  "Hi World",
		},
	}

	// for i := 0; i < 1000000; i++ {
	//for i := 0; i < 100000; i++ {
	for i := 0; i < 10000; i++ {
		n := fmt.Sprintf("name-%d", i)
		c.Create[n] = data.Object{
			"_id":   n,
			"_type": "location",
			"name":  n,
		}
	}
	//*/

	fmt.Printf("=======================================\n")

	data.Update(c)

	wait_for_TERM_or_INT()

	data.Stop()
}

func wait_for_TERM_or_INT() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}

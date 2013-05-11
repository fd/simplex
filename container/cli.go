package container

import (
	"fmt"
	"github.com/fd/options"
	"os"
	"path"
	"strings"
)

const opts_spec_tmpl = `
useage: {{CMD}}
Semi-static web framework
--
port=3000 -p,--port,PORT    HTTP server port
addr=     -a,--addr,ADDR    HTTP server address (use instead of --port)
src=      --src,DATA_SRC    Source data store
dst=      --dst,DATA_DST    Source data store
--
--
server    server,s          Run the web server
generate  generate,g        Generate the aggregate data
--
`

func CLI() {
	var (
		opts_spec_str = strings.Replace(opts_spec_tmpl, "{{CMD}}", path.Base(os.Args[0]), -1)
		opts_spec     = options.MustParse(opts_spec_str)
		opts          = opts_spec.MustInterpret(os.Args, os.Environ())
		env           Environment
	)

	if s := opts.Get("addr"); s != "" {
		env.HttpAddr = s
	}

	if s := opts.Get("port"); s != "" {
		env.HttpAddr = ":" + s
	}

	if s := opts.Get("src"); s != "" {
		env.Source = s
	}

	if s := opts.Get("dst"); s != "" {
		env.Destination = s
	}

	switch opts.Command {

	case "server":
		err := Serve(env)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	case "generate":
		err := Generate(env)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	default:
		opts_spec.PrintUsageAndExit()

	}
}

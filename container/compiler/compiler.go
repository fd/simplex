package compiler

import (
	"fmt"
	template "github.com/fd/w/template/compiler"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
)

const src = `package main

// Libraries
import (
  "github.com/fd/w/container"
)

// Apps
import (
%s
)

func main() {
  container.Run()
}
`

func WriteContainerWrapper(ctx *template.Context) error {
	wrapper := PrintContainerWrapper(ctx)

	filename := path.Join(ctx.WROOT, "wrapper.go")

	err := ioutil.WriteFile(filename, []byte(wrapper), 0644)
	return err
}

func RemoveContainerWrapper(ctx *template.Context) error {
	filename := path.Join(ctx.WROOT, "wrapper.go")

	err := os.Remove(filename)
	return err
}

func PrintContainerWrapper(ctx *template.Context) string {
	apps := []string{}
	for app, _ := range ctx.Applications {
		apps = append(apps, fmt.Sprintf("  _ %s", strconv.Quote(app)))
	}

	wrapper := fmt.Sprintf(src, strings.Join(apps, "\n"))

	return wrapper
}

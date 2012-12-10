package util

import (
	"path"
	"runtime"
	"strings"
)

func InitializingPackage() string {
	pc := make([]uintptr, 1024)
	n := runtime.Callers(0, pc)

	for i := 0; i < n; i++ {
		f := runtime.FuncForPC(pc[i])
		if f == nil {
			continue
		}

		name := f.Name()

		if !strings.HasSuffix(name, ".init") {
			continue
		}

		file, _ := f.FileLine(pc[i])
		pkg_path := path.Dir(file)

		return pkg_path
	}

	return "unknown"
}

func InitializingApplication() string {
	pkg := InitializingPackage()
	parts := strings.Split("/")

	for i, part := range parts {
		if part == "apps" {
			if len(parts) > (i + 1) {
				return strings.Join(parts[:i+1], "/")
			} else {
				break
			}
		}
	}

	return "unknown"
}

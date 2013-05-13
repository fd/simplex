package stesting

import (
	"os"
	"path"
	"runtime"
	"simplex.sh/static"
	"simplex.sh/store"
	_ "simplex.sh/store/file"
	"testing"
)

func Golden(t *testing.T, g static.Generator) {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	src, err := store.Open("file://" + path.Join(wd, "test/src"))
	if err != nil {
		t.Fatal(err)
	}

	dst, err := store.Open("file://" + path.Join(wd, "test/dst"))
	if err != nil {
		t.Fatal(err)
	}

	os.RemoveAll(path.Join(wd, "test/dst"))

	err = static.Generate(
		src,
		dst,
		g,
	)

	if err != nil {
		t.Fatal(err)
	}
}

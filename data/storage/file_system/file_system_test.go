package file_system

import (
	storageT "github.com/fd/simplex/data/storage/testing"
	"io/ioutil"
	"os"
	"testing"
)

func TestFileSystem(t *testing.T) {
	root, err := ioutil.TempDir("", "w-testing-")
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		os.RemoveAll(root)
	}()

	storageT.ValidateDriver(t, &S{
		Root: root,
	})
}

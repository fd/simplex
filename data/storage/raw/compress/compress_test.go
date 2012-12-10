package compress

import (
	"github.com/fd/w/data/storage/raw/file_system"
	storageT "github.com/fd/w/data/storage/testing"
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

	storageT.ValidateRawDriver(t, &S{Storage: &file_system.S{
		Root: root,
	}})
}

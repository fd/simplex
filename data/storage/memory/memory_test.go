package memory

import (
	storageT "github.com/fd/w/data/storage/testing"
	"testing"
)

func TestFileSystem(t *testing.T) {
	storageT.ValidateDriver(t, &S{
		objects: map[string][]byte{},
	})
}

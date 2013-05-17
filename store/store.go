package store

import (
	"io"
)

type Store interface {
	GetBlob(name string) (io.ReadCloser, error)
	SetBlob(name string) (io.WriteCloser, error)
}

type notfound_error struct {
	path string
}

func IsNotFound(err error) bool {
	_, ok := err.(*notfound_error)
	return ok
}

func NotFoundError(path string) error {
	return &notfound_error{path}
}

func (n *notfound_error) Error() string {
	return "Object not found: " + n.path
}

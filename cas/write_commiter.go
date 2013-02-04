package cas

import (
	"io"
)

type Commiter interface {
	Commit(addr Addr) error
	Rollback() error
}

type WriteCommiter interface {
	io.Writer
	Commiter
}

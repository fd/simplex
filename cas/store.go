package cas

import (
	"io"
)

type Getter interface {
	Get(Addr) (io.ReadCloser, error)
}

type Setter interface {
	Set() (WriteCommiter, error)
}

type GetterSetter interface {
	Getter
	Setter
}

type Closer interface {
	Close() error
}

type Store interface {
	GetterSetter
	Closer
}

package cas

import (
	"bytes"
	"encoding/hex"
)

type Addr []byte

type addr_kind byte

const (
	addr_kind__uncompressed_val addr_kind = iota
	addr_kind__compressed_val
	addr_kind__sha
	addr_kind__ref
)

func (a Addr) String() string {
	return hex.EncodeToString(a)
}

func CompareAddr(a, b Addr) int {
	return bytes.Compare(a, b)
}

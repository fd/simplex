package cas

import (
	"bytes"
	"encoding/hex"
)

type Addr []byte

func ParseAddr(str string) (Addr, error) {
	return hex.DecodeString(str)
}

func (a Addr) String() string {
	return hex.EncodeToString(a)
}

func (a Addr) Compare(b Addr) int {
	return bytes.Compare(a, b)
}

package ident

import (
	"crypto/sha1"
	"encoding/hex"
)

type SHA [20]byte

var ZeroSHA = SHA{}

func Hash(val interface{}) SHA {
	return HashCompairBytes(CompairBytes(val))
}

func HashCompairBytes(dat []byte) SHA {
	h := sha1.New()
	h.Write(dat)
	b := h.Sum(nil)
	s := SHA{}
	copy(s[:], b)
	return s
}

func HexSHA(sha SHA) string {
	return hex.EncodeToString(sha[:])
}

package data

import (
	"encoding/gob"
)

type Value interface{}
type Object map[string]interface{}

func init() {
	gob.Register(Object{})
}

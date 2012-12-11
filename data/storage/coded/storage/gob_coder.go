package storage

import (
	"bytes"
	"encoding/gob"
)

type GobCoder struct {
}

func (c *GobCoder) Encode(v interface{}) ([]byte, error) {
	var b bytes.Buffer
	err := gob.NewEncoder(&b).Encode(v)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (c *GobCoder) Decode(dat []byte, val interface{}) error {
	err := gob.NewDecoder(bytes.NewBuffer(dat)).Decode(val)
	if err != nil {
		return err
	}
	return nil
}

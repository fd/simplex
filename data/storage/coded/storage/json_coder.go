package storage

import (
	"encoding/json"
)

type JsonCoder struct {
}

func (c *JsonCoder) Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (c *JsonCoder) Decode(dat []byte, val interface{}) error {
	err := json.Unmarshal(dat, val)
	if err != nil {
		return err
	}
	return nil
}

package storage

import (
	"encoding/json"
)

type JsonCoder struct {
}

func (c *JsonCoder) Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (c *JsonCoder) Decode(dat []byte) (interface{}, error) {
	var val interface{}
	err := json.Unmarshal(dat, &val)
	if err != nil {
		return nil, err
	}
	return val, nil
}

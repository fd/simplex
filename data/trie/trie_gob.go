package trie

import (
	"bytes"
	"compress/zlib"
	"encoding/gob"
)

func (t *T) GobDecode(data []byte) error {
	buf := bytes.NewReader(data)

	comp, err := zlib.NewReader(buf)
	if err != nil {
		return err
	}
	defer comp.Close()

	err = gob.NewDecoder(comp).Decode(&t.root)

	if t.root != nil {
		t.root.set_parents()
	}

	return err
}

func (t *T) GobEncode() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 1024*1024))

	comp, err := zlib.NewWriterLevel(buf, zlib.DefaultCompression)
	if err != nil {
		return nil, err
	}
	defer comp.Close()

	err = gob.NewEncoder(comp).Encode(t.root)

	comp.Close()

	return buf.Bytes(), err
}

func (n *node_t) set_parents() {
	for _, c := range n.Children {
		c.parent = n
		c.set_parents()
	}
}

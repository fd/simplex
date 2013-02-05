package cas

import (
	"bytes"
	"compress/zlib"
	"github.com/fd/simplex/cas/blob"
	"io"
	"io/ioutil"
	"reflect"
)

type Decoder struct {
	err error

	outbound_r io.ReadCloser

	blob_dec *blob.Decoder
}

func Decode(s GetterSetter, addr Addr, e interface{}) error {
	dec := NewDecoder(s, addr)

	err := dec.Decode(e)
	if err != nil {
		return err
	}

	err = dec.Close()
	if err != nil {
		return err
	}

	return nil
}

func NewDecoder(s GetterSetter, addr Addr) *Decoder {
	var (
		err error
		dec = &Decoder{}
	)

	switch addr_kind(addr[0]) {

	case addr_kind__uncompressed_val:
		dec.outbound_r = ioutil.NopCloser(bytes.NewReader(addr[1:]))

	case addr_kind__compressed_val:
		r := ioutil.NopCloser(bytes.NewReader(addr[1:]))

		r, err = zlib.NewReader(r)
		if err != nil {
			dec.err = err
			return dec
		}

		dec.outbound_r = r

	case addr_kind__sha:
		r, err := s.Get(addr)
		if err != nil {
			dec.err = err
			return dec
		}

		r, err = zlib.NewReader(r)
		if err != nil {
			dec.err = err
			return dec
		}

		dec.outbound_r = r

	default:
		panic("not reached")

	}

	return dec
}

func (dec *Decoder) Read(p []byte) (n int, err error) {
	if dec.err != nil {
		return 0, err
	}

	return dec.outbound_r.Read(p)
}

func (dec *Decoder) Close() error {
	if dec.err != nil {
		return dec.err
	}

	return dec.outbound_r.Close()
}

func (dec *Decoder) Decode(e interface{}) error {
	if dec.blob_dec == nil {
		dec.blob_dec = blob.NewDecoder(dec.outbound_r)
	}

	return dec.blob_dec.Decode(e)
}

func (dec *Decoder) DecodeValue(e reflect.Value) error {
	if dec.blob_dec == nil {
		dec.blob_dec = blob.NewDecoder(dec.outbound_r)
	}

	return dec.blob_dec.DecodeValue(e)
}

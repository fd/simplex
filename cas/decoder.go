package cas

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"simplex.sh/cas/blob"
)

type Decoder struct {
	err error

	inbound_r  io.Closer
	outbound_r io.ReadCloser

	blob_dec *blob.Decoder
}

func Decode(s GetterSetter, addr Addr, e interface{}) error {
	dec := NewDecoder(s, addr)
	defer dec.Close()

	err := dec.Decode(e)
	if err != nil {
		return err
	}

	return nil
}

func DecodeValue(s GetterSetter, addr Addr, e reflect.Value) error {
	dec := NewDecoder(s, addr)
	defer dec.Close()

	err := dec.DecodeValue(e)
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

	if len(addr) < 1 {
		dec.err = fmt.Errorf("Cannot decode empty address.")
		return dec
	}

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
		dec.inbound_r = r

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
	if dec.inbound_r != nil {
		dec.inbound_r.Close()
	}

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

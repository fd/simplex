package cas

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"github.com/fd/simplex/cas/blob"
	"hash"
	"io"
	"reflect"
)

const DEFAULT_OVERFLOW_TRIGGER = 256

type Encoder struct {
	Addr Addr

	err              error
	overflow_trigger int

	inbound_w io.Writer

	outbound_c   Commiter
	compressed_c io.Closer

	uncompressed_b *bytes.Buffer
	compressed_b   *bytes.Buffer

	hash hash.Hash

	blob_enc *blob.Encoder
}

func Encode(s GetterSetter, e interface{}, overflow_trigger int) (Addr, error) {
	enc := NewEncoder(s, overflow_trigger)

	err := enc.Encode(e)
	if err != nil {
		return nil, err
	}

	err = enc.Close()
	if err != nil {
		return nil, err
	}

	return enc.Addr, nil
}

func NewEncoder(store GetterSetter, overflow_trigger int) *Encoder {
	if overflow_trigger < 0 {
		overflow_trigger = DEFAULT_OVERFLOW_TRIGGER
	}

	outbound_w, err := store.Set()
	if err != nil {
		return &Encoder{err: err}
	}

	hash_w := sha1.New()

	pre_hash_writer := io.MultiWriter(hash_w, outbound_w)

	compressed_b := bytes.NewBuffer(make([]byte, 0, overflow_trigger))
	compressed_w, err := zlib.NewWriterLevel(&overflow_writer{
		trigger:     overflow_trigger,
		overflow_w:  pre_hash_writer,
		underflow_w: compressed_b,
	}, zlib.DefaultCompression)

	if err != nil {
		return &Encoder{err: err}
	}

	uncompressed_b := bytes.NewBuffer(make([]byte, 0, overflow_trigger))
	uncompressed_w := &overflow_writer{
		trigger:     overflow_trigger,
		overflow_w:  compressed_w,
		underflow_w: uncompressed_b,
	}

	return &Encoder{
		overflow_trigger: overflow_trigger,
		inbound_w:        uncompressed_w,
		outbound_c:       outbound_w,
		compressed_c:     compressed_w,
		uncompressed_b:   uncompressed_b,
		compressed_b:     compressed_b,
		hash:             hash_w,
	}
}

func (enc *Encoder) Write(p []byte) (n int, err error) {
	if enc.err != nil {
		return 0, enc.err
	}
	return enc.inbound_w.Write(p)
}

func (enc *Encoder) Close() error {
	if enc.err != nil {
		return enc.err
	}

	err := enc.compressed_c.Close()
	if err != nil {
		enc.outbound_c.Rollback()
		return err
	}

	if enc.overflow_trigger > 0 && enc.uncompressed_b.Len() < enc.overflow_trigger {
		b := make([]byte, 1+enc.uncompressed_b.Len())
		b[0] = byte(addr_kind__uncompressed_val)
		copy(b[1:], enc.uncompressed_b.Bytes())
		enc.Addr = Addr(b)
		enc.outbound_c.Rollback()
		return nil
	}

	if enc.overflow_trigger > 0 && enc.compressed_b.Len() < enc.overflow_trigger {
		b := make([]byte, 1+enc.compressed_b.Len())
		b[0] = byte(addr_kind__compressed_val)
		copy(b[1:], enc.compressed_b.Bytes())
		enc.Addr = Addr(b)
		enc.outbound_c.Rollback()
		return nil
	}

	sum := enc.hash.Sum(nil)
	b := make([]byte, 1+len(sum))
	b[0] = byte(addr_kind__sha)
	copy(b[1:], sum)
	enc.Addr = Addr(b)

	err = enc.outbound_c.Commit(enc.Addr)
	if err != nil {
		return err
	}

	return nil
}

func (enc *Encoder) Encode(e interface{}) error {
	if enc.blob_enc == nil {
		enc.blob_enc = blob.NewEncoder(enc.inbound_w)
	}

	return enc.blob_enc.Encode(e)
}

func (enc *Encoder) EncodeValue(e reflect.Value) error {
	if enc.blob_enc == nil {
		enc.blob_enc = blob.NewEncoder(enc.inbound_w)
	}

	return enc.blob_enc.EncodeValue(e)
}

type overflow_writer struct {
	trigger     int
	n           int
	overflow_w  io.Writer
	underflow_w *bytes.Buffer
}

func (w *overflow_writer) Write(p []byte) (n int, err error) {
	if w.n >= w.trigger {
		return w.overflow_w.Write(p)
	}

	n, err = w.underflow_w.Write(p)
	w.n += n
	if err != nil {
		return
	}

	if w.n >= w.trigger {
		return w.overflow_w.Write(w.underflow_w.Bytes())
	}

	return
}

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

const (
	st_internal encoder_state = iota
	st_internal_compressed
	st_external
)

type (
	encoder_state int

	Encoder struct {
		// Addr is non nil after the encoder is closed and no errors
		// occured.
		Addr Addr

		overflow_trigger int
		store            Setter

		err        error
		state      encoder_state
		log_buffer *bytes.Buffer
		writer     io.Writer

		uncompressed_b *bytes.Buffer
		compressed_w   *zlib.Writer
		compressed_b   *bytes.Buffer
		hash           hash.Hash
		outbound_w     WriteCommiter

		blob_enc *blob.Encoder
	}
)

func Encode(s Setter, e interface{}, overflow_trigger int) (Addr, error) {
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

func NewEncoder(store Setter, overflow_trigger int) *Encoder {
	if overflow_trigger < 0 {
		overflow_trigger = DEFAULT_OVERFLOW_TRIGGER
	}

	log_buffer := bytes.NewBuffer(nil)
	uncompressed_b := bytes.NewBuffer(nil)
	writer := io.MultiWriter(log_buffer, uncompressed_b)

	return &Encoder{
		store:            store,
		overflow_trigger: overflow_trigger,
		uncompressed_b:   uncompressed_b,
		writer:           writer,
		log_buffer:       log_buffer,
	}
}

func (enc *Encoder) Write(p []byte) (n int, err error) {
	if enc.err != nil {
		return 0, enc.err
	}

	switch enc.state {

	case st_internal:
		n, err := enc.writer.Write(p)
		if err != nil {
			enc.err = err
		}
		if enc.uncompressed_b.Len() > enc.overflow_trigger {
			enc.switch_to(st_internal_compressed)
		}
		return n, enc.err

	case st_internal_compressed:
		n, err := enc.writer.Write(p)
		if err != nil {
			enc.err = err
		}
		if enc.compressed_b.Len() > enc.overflow_trigger {
			enc.switch_to(st_external)
		}
		return n, enc.err

	case st_external:
		n, err := enc.writer.Write(p)
		if err != nil {
			enc.err = err
		}
		return n, enc.err

	}

	panic("not reached")
}

func make_addr(kind addr_kind, data []byte) Addr {
	b := make([]byte, 1+len(data))
	b[0] = byte(kind)
	copy(b[1:], data)
	return Addr(b)
}

func (enc *Encoder) Close() error {
	if enc.err != nil {
		return enc.err
	}

	switch enc.state {

	case st_internal:
		enc.Addr = make_addr(
			addr_kind__uncompressed_val,
			enc.uncompressed_b.Bytes(),
		)
		return nil

	case st_internal_compressed:
		err := enc.compressed_w.Close()
		if err != nil {
			enc.err = err
			return enc.err
		}
		enc.Addr = make_addr(
			addr_kind__compressed_val,
			enc.compressed_b.Bytes(),
		)
		return nil

	case st_external:
		err := enc.compressed_w.Close()
		if err != nil {
			enc.outbound_w.Rollback()
			enc.err = err
			return enc.err
		}

		addr := make_addr(
			addr_kind__sha,
			enc.hash.Sum(nil),
		)

		err = enc.outbound_w.Commit(addr)
		if err != nil {
			enc.err = err
			return enc.err
		}

		enc.Addr = addr

		return nil

	}

	panic("not reached")
}

func (enc *Encoder) Encode(e interface{}) error {
	if enc.blob_enc == nil {
		enc.blob_enc = blob.NewEncoder(enc)
	}

	return enc.blob_enc.Encode(e)
}

func (enc *Encoder) EncodeValue(e reflect.Value) error {
	if enc.blob_enc == nil {
		enc.blob_enc = blob.NewEncoder(enc)
	}

	return enc.blob_enc.EncodeValue(e)
}

func (enc *Encoder) switch_to(state encoder_state) {
	switch state {

	case st_internal_compressed:
		// setup
		compressed_b := bytes.NewBuffer(nil)
		compressed_w, err := zlib.NewWriterLevel(compressed_b, zlib.DefaultCompression)
		if err != nil {
			enc.err = err
			return
		}

		// write log buffer
		_, err = compressed_w.Write(enc.log_buffer.Bytes())
		if err != nil {
			enc.err = err
			return
		}

		writer := io.MultiWriter(enc.log_buffer, compressed_w)
		enc.writer = writer
		enc.uncompressed_b = nil
		enc.compressed_b = compressed_b
		enc.compressed_w = compressed_w

		err = enc.compressed_w.Flush()
		if err != nil {
			enc.err = err
			return
		}

		enc.state = st_internal_compressed
		if enc.compressed_b.Len() > enc.overflow_trigger {
			enc.switch_to(st_external)
		}

		return

	case st_external:
		// setup
		hash := sha1.New()
		outbound_w, err := enc.store.Set()
		if err != nil {
			enc.err = err
			return
		}

		compressed_w, err := zlib.NewWriterLevel(
			io.MultiWriter(hash, outbound_w),
			zlib.DefaultCompression,
		)
		if err != nil {
			enc.err = err
			return
		}

		// write log buffer
		_, err = compressed_w.Write(enc.log_buffer.Bytes())
		if err != nil {
			enc.err = err
			return
		}

		writer := io.MultiWriter(enc.log_buffer, compressed_w)
		enc.writer = writer
		enc.uncompressed_b = nil
		enc.compressed_b = nil
		enc.compressed_w = compressed_w
		enc.outbound_w = outbound_w
		enc.hash = hash

		err = enc.compressed_w.Flush()
		if err != nil {
			enc.err = err
			return
		}

		enc.state = st_external
		return

	}

	panic("not reached")
}

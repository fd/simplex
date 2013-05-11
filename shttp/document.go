package shttp

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"hash"
	"io"
	"net/http"
	"strconv"
)

type document struct {
	Digest string
	Status int
	Header http.Header
	Body   []byte
}

type Writer interface {
	http.ResponseWriter
	Router
}

type document_writer struct {
	io.Writer
	route_builder
	digest   hash.Hash
	body     bytes.Buffer
	document *document
}

func new_document_writer() *document_writer {
	w := &document_writer{
		digest: md5.New(),
		document: &document{
			Header: make(http.Header, 10),
		},
	}

	w.Writer = io.MultiWriter(&w.body, w.digest)

	return w
}

func (d *document_writer) Header() http.Header {
	return d.document.Header
}

func (d *document_writer) Close() error {
	d.document.Body = d.body.Bytes()
	d.document.Digest = hex.EncodeToString(d.digest.Sum(nil))

	if d.document.Status == 0 {
		d.document.Status = 200
	}

	if d.document.Header.Get("Content-Type") == "" {
		d.document.Header.Set("Content-Type", "text/html; charset=utf-8")
	}

	d.document.Header.Set("Content-Length", strconv.Itoa(len(d.document.Body)))
	d.document.Header.Set("ETag", strconv.Quote(d.document.Digest))

	return nil
}

func (d *document_writer) WriteHeader(status int) {
	d.document.Status = status
}

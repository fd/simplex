package shttp

import (
	"database/sql"
	"encoding/hex"
	"net/http"
	"simplex.sh/store/cas"
	"strconv"
)

type document struct {
	Digest string
	Status int
	Header http.Header
}

type Writer interface {
	http.ResponseWriter
	Router
}

type document_writer struct {
	cas.Writer
	route_builder
	document *document
}

func new_document_writer(tx *sql.Tx) *document_writer {
	w := &document_writer{
		Writer: cas.OpenWriter(tx),
		document: &document{
			Header: make(http.Header, 10),
		},
	}

	return w
}

func (d *document_writer) Header() http.Header {
	return d.document.Header
}

func (d *document_writer) Close() error {
	err := d.Writer.Close()
	if err != nil {
		return err
	}

	addr := d.Address()
	d.document.Digest = hex.EncodeToString(addr[:])

	if d.document.Status == 0 {
		d.document.Status = 200
	}

	if d.document.Header.Get("Content-Type") == "" {
		d.document.Header.Set("Content-Type", "text/html; charset=utf-8")
	}

	d.document.Header.Set("Content-Length", strconv.Itoa(d.Len()))
	d.document.Header.Set("ETag", strconv.Quote(d.document.Digest))

	return nil
}

func (d *document_writer) WriteHeader(status int) {
	d.document.Status = status
}

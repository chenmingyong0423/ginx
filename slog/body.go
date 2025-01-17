package slog

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	_ http.ResponseWriter = (*respBodyWriter)(nil)
	_ http.Flusher        = (*respBodyWriter)(nil)
	_ http.Hijacker       = (*respBodyWriter)(nil)
)

type reqBodyReader struct {
	io.ReadCloser
	body  *bytes.Buffer
	bytes int64
}

func (r *reqBodyReader) Read(b []byte) (n int, err error) {
	n, err = r.ReadCloser.Read(b)
	r.bytes += int64(n)
	if r.body != nil {
		r.body.Write(b[:n])
	}
	return
}

type respBodyWriter struct {
	gin.ResponseWriter
	body  *bytes.Buffer
	bytes int64
}

func (w *respBodyWriter) Write(b []byte) (n int, err error) {
	n, err = w.ResponseWriter.Write(b)
	w.bytes += int64(n)
	if w.body != nil {
		w.body.Write(b[:n])
	}
	return
}

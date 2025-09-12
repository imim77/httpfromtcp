package response

import (
	"fmt"
	"io"

	"github.com/imim77/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	StatusOK          StatusCode = 200
	StatusBadRequest  StatusCode = 400
	StatusServerError StatusCode = 500
)

func GetDefaultHeaders(contentLen int) *headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/html")
	return h
}

type Writer struct {
	writer io.Writer
}

func NewWriter(writer io.Writer) *Writer {
	return &Writer{
		writer: writer,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	mapa := map[StatusCode]string{
		StatusOK:          "OK",
		StatusBadRequest:  "Bad Request",
		StatusServerError: "Internal Server Error",
	}

	statusLine := fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, mapa[statusCode])

	_, err := w.writer.Write([]byte(statusLine))
	if err != nil {
		return fmt.Errorf("unrecognized error code")
	}
	return nil
}
func (w *Writer) WriteHeaders(h headers.Headers) error {
	var err error = nil
	b := []byte{}
	h.ForEach(func(n, v string) {
		b = fmt.Append(b, fmt.Sprintf("%s: %s\r\n", n, v))

	})
	b = fmt.Append(b, "\r\n")
	_, err = w.writer.Write(b)
	if err != nil {
		return err
	}
	return err

}
func (w *Writer) WriteBody(p []byte) (int, error) {
	n, err := w.writer.Write(p)
	if err != nil {
		return 0, err
	}
	return n, nil

}

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

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	mapa := map[StatusCode]string{
		StatusOK:          "OK",
		StatusBadRequest:  "Bad Request",
		StatusServerError: "Internal Server Error",
	}

	statusLine := fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, mapa[statusCode])

	_, err := w.Write([]byte(statusLine))
	if err != nil {
		return fmt.Errorf("unrecognized error code")
	}

	return nil
}

func GetDefaultHeaders(contentLen int) *headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")
	return h
}

func WriteHeaders(w io.Writer, headers *headers.Headers) error {
	var err error = nil
	b := []byte{}
	headers.ForEach(func(n, v string) {
		b = fmt.Append(b, fmt.Sprintf("%s: %s\r\n", n, v))

	})
	b = fmt.Append(b, "\r\n")
	_, err = w.Write(b)
	if err != nil {
		return err
	}
	return err
}

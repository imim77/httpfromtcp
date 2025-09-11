package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/imim77/httpfromtcp/internal/headers"
)

type praserState string

const (
	StateInit    praserState = "init"
	StateDone    praserState = "done"
	StateHeaders praserState = "headers"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	Headers     *headers.Headers
	RequestLine RequestLine
	state       praserState
}

func newRequest() *Request {
	return &Request{
		state:   StateInit,
		Headers: headers.NewHeaders(),
	}
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0
outer:
	for {
		currentdata := data[read:]

		switch r.state {
		case StateInit:

			rl, n, err := parseRequestLine(currentdata)
			if err != nil {
				return 0, err
			}
			if n == 0 {
				break outer
			}
			r.RequestLine = *rl
			read += n
			r.state = StateHeaders
		case StateDone:
			break outer
		case StateHeaders:

			n, done, err := r.Headers.Parse(currentdata)
			if err != nil {
				return 0, err
			}

			if n == 0 { // we couldn't read anything
				break outer
			}

			read += n

			if done {
				r.state = StateDone
			}
		default:
			panic("somehow shit")
		}

	}
	return read, nil
}

func (r *Request) done() bool {
	return r.state == StateDone
}

var SEPERATOR = []byte("\r\n")
var SEP = []byte("/")

func parseRequestLine(b []byte) (*RequestLine, int, error) {
	ind := bytes.Index(b, SEPERATOR)
	if ind == -1 {
		return nil, 0, nil
	}

	startLine := b[:ind] // npr GET / HTTP / 1.1 \r\n -- kraj ..........
	read := ind + len(SEPERATOR)

	parts := bytes.Fields(startLine)
	if len(parts) != 3 { // if we don't have method,version, http proctocol
		return nil, 0, errors.Join(fmt.Errorf("don't have all the parts of the request line"))
	}

	httpPart := bytes.Split(parts[2], SEP)
	if len(httpPart) != 2 || string(httpPart[0]) != "HTTP" || string(httpPart[1]) != "1.1" {
		return nil, 0, errors.Join(fmt.Errorf("http version not valid"))
	}

	rl := &RequestLine{
		Method:        string(parts[0]),
		RequestTarget: string(parts[1]),
		HttpVersion:   string(httpPart[1]),
	}
	return rl, read, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()
	buf := make([]byte, 1024)
	bufLen := 0
	for !request.done() {
		n, err := reader.Read(buf[bufLen:])
		if err != nil {
			return nil, errors.Join(fmt.Errorf("error while reading"), err)
		}
		bufLen += n
		readN, err := request.parse(buf[:bufLen])
		if err != nil {
			return nil, err
		}
		copy(buf, buf[readN:bufLen])
		bufLen -= readN

	}
	return request, nil
}

package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/imim77/httpfromtcp/internal/headers"
)

type praserState string

const (
	StateInit    praserState = "init"
	StateDone    praserState = "done"
	StateHeaders praserState = "headers"
	StateBody    praserState = "body"
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
	Body        []byte
}

func getInt(headers *headers.Headers, name string, defaultValue int) int {
	valueStr, exists := headers.Get(name)
	if !exists {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0
	}
	return value

}

func newRequest() *Request {
	return &Request{
		state:   StateInit,
		Headers: headers.NewHeaders(),
		Body:    []byte(""),
	}
}

func (r *Request) hasBody() bool {
	length := getInt(r.Headers, "content-length", 0) // total size of the body according to the "content-length" header
	return length > 0
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0
outer:
	for {
		currentdata := data[read:]
		if len(currentdata) == 0 {
			break outer
		}

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
				if r.hasBody() {
					r.state = StateBody
				} else {
					r.state = StateDone
				}
			}
		case StateBody:
			length := getInt(r.Headers, "content-length", 0) // total size of the body according to the "content-length" header
			if length == 0 {
				panic("chunked not implemented yet")
			}
			// r.Body => already read part of the body
			// ammount we have left to read or the ammount of the data(for example 20B) that is avaliable in the current block
			remaining := min(length-len(r.Body), len(currentdata))
			r.Body = append(r.Body, currentdata[:remaining]...)
			read += remaining

			if len(r.Body) == length {
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

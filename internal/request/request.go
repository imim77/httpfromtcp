package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
}

var SEPERATOR = "\r\n"

func parseRequestLine(b string) (*RequestLine, string, error) {
	ind := strings.Index(b, SEPERATOR)
	if ind == -1 {
		return nil, b, nil
	}

	startLine := b[:ind]                // npr GET / HTTP / 1.1 \r\n -- kraj ..........
	restOfMsq := b[ind+len(SEPERATOR):] // field lines

	parts := strings.Fields(startLine)
	if len(parts) != 3 { // if we don't have method,version, http proctocol
		return nil, b, errors.Join(fmt.Errorf("don't have all the parts of the request line"))
	}

	httpPart := strings.Split(parts[2], "/")
	if len(httpPart) != 2 || httpPart[0] != "HTTP" || httpPart[1] != "1.1" {
		return nil, b, errors.Join(fmt.Errorf("http version not valid"))
	}

	rl := &RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   httpPart[1],
	}
	return rl, restOfMsq, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("unable to io.ReadAll"), err)
	}
	rl, _, err := parseRequestLine(string(data))
	if err != nil {
		return nil, errors.Join(fmt.Errorf("unable to parse request line"), err)
	}
	return &Request{RequestLine: *rl}, nil
}

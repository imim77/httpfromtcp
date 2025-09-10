package headers

import (
	"bytes"
	"fmt"
)

type Headers map[string]string

var rn = []byte("\r\n")

func NewHeaders() Headers {
	return map[string]string{}
}

func parseHeader(fieldLine []byte) (string, string, error) { // "Host: example.com\r\n"
	parts := bytes.SplitN(fieldLine, []byte(":"), 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("bad header")
	}
	name := parts[0]                   // "Host"
	value := bytes.TrimSpace(parts[1]) // "example.com"

	if bytes.HasSuffix(name, []byte(" ")) {
		return "", "", fmt.Errorf("bad header")
	}
	return string(name), string(value), nil

}

func (h Headers) Parse(data []byte) (int, bool, error) {
	read := 0
	done := false
	for {
		ind := bytes.Index(data[read:], rn) //  the length of the string without \r\n
		fmt.Printf("parsing header(%d) - %d\n", read, ind)
		if ind == -1 {
			break

		}
		// WE HAVE HIT EMPTY HEADER
		if ind == 0 { // if the \r\n is on the first index, if there is not anything infront of it, that is represeting the end of the header
			done = true
			read += len(rn)
			break
		}

		name, value, err := parseHeader(data[read : read+ind]) // from 0 to 17()
		if err != nil {
			return 0, false, err
		}
		read += ind + len(rn) // len(rn) for "jumping across" the rn value which is one index i
		h[name] = value

	}
	return read, done, nil

}

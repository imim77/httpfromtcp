package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers struct {
	headers map[string]string
}

var rn = []byte("\r\n")

func NewHeaders() *Headers {
	return &Headers{
		headers: map[string]string{},
	}
}

func (h *Headers) Set(key, val string) {
	key = strings.ToLower(key)

	if v, ok := h.headers[key]; ok {
		h.headers[key] = fmt.Sprintf("%s,%s", v, val)
	} else {
		h.headers[key] = val
	}
}

func (h *Headers) Get(key string) string {
	return h.headers[strings.ToLower(key)]
}

func (h *Headers) ForEach(cb func(k, v string)) {
	for k, v := range h.headers {
		cb(k, v)
	}
}

func parseHeader(filedLine []byte) (string, string, error) {
	parts := bytes.SplitN(filedLine, []byte(":"), 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("malformed field line")
	}

	key := parts[0]
	value := bytes.TrimSpace(parts[1])

	if bytes.HasSuffix(key, []byte(" ")) {
		return "", "", fmt.Errorf("malformed filed name")
	}

	return string(key), string(value), nil

}

func (h *Headers) Parse(data []byte) (int, bool, error) {
	read := 0
	done := false
	for {
		idx := bytes.Index(data[read:], rn) //  the length of the string without \r\n
		if idx == -1 {
			break
		}
		//EMPTY HEADER
		if idx == 0 { // if the \r\n is on the first index, if there is not anything infront of it, that is represeting the end of the header
			done = true
			read += len(rn)
			break
		}
		key, val, err := parseHeader(data[read : read+idx])
		if err != nil {
			return 0, false, fmt.Errorf("error parsing the data")
		}
		if !isToken([]byte(key)) {
			return 0, false, fmt.Errorf("invalid field-name(token)")
		}
		read += idx + len(rn) // len(rn) for "jumping across" the rn value which is one index i
		h.Set(key, val)

	}
	return read, done, nil
}

//  field-name     = token (RFC 9110)
// MUST CONTAIN
//Uppercase letters: A-Z
//Lowercase letters: a-z
//Digits: 0-9
//Special characters: !, #, $, %, &, ', *, +, -, ., ^, _, `, |, ~

func isToken(data []byte) bool {
	for _, ch := range data {
		found := false
		if ch >= 'A' && ch <= 'Z' || ch >= 'a' && ch <= 'z' || ch >= '0' && ch <= '9' {
			found = true
		}
		switch ch {
		case '!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~':
			found = true
		}
		if !found {
			return false
		}
	}
	return true
}

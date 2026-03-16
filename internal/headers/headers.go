package headers

import (
	"bytes"
	"fmt"
	"strings"
)

const (
	CRLF = "\r\n"
)

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{}
}
func (h Headers) Parse(data []byte) (n int, done bool, err error) {

	i := bytes.Index(data, []byte(CRLF))
	if i <= 0 {
		return 0, i == 0, nil
	}
	ln := string(data[:i])
	data = data[i+len(CRLF):]

	key, val, _ := strings.Cut(ln, ":")
	if key != strings.TrimSpace(key) {
		return 0, false, fmt.Errorf("invalid header name")
	}

	val = strings.TrimSpace(val)
	h[key] = val
	return i, true, nil
}

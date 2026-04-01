package headers

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

const (
	CRLF = "\r\n"
)

type Headers struct {
	headers map[string]string
}

func NewHeaders() *Headers {
	return &Headers{
		headers: map[string]string{},
	}
}

func (h *Headers) HeadersMap() map[string]string {
	return h.headers
}

func (h *Headers) Get(key string) string {
	return h.headers[strings.ToLower(key)]
}

func (h *Headers) Set(key string, val string) {
	if val, ok := h.headers[key]; ok {
		h.headers[key] = h.headers[key] + ", " + val
		return
	}
	h.headers[strings.ToLower(key)] = val
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {

	i := bytes.Index(data, []byte(CRLF))
	if i <= 0 {
		n := 0
		if i == 0 {
			n = len(CRLF)
		}
		return n, i == 0, nil
	}
	ln := string(data[:i])
	data = data[i+len(CRLF):]

	key, val, _ := strings.Cut(ln, ":")
	match, err := regexp.MatchString("^[a-zA-Z!#$%&'*+-.^_`|~]+$", key)

	if err != nil {
		return 0, false, fmt.Errorf("error parsing headers")
	}

	if key != strings.TrimSpace(key) || !match {
		return 0, false, fmt.Errorf("invalid header name")
	}
	val = strings.TrimSpace(val)

	h.Set(key, val)
	return i + len(CRLF), false, nil
}

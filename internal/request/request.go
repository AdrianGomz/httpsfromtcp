package request

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
)

const (
	SEPARATOR = "\r\n"
	HTTP_V1   = "HTTP/1.1"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	Headers     map[string]string
	Body        []byte
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("unable to io.ReadAll: %w", err)
	}
	rl, err := parseRequestLine(strings.Split(string(data), "\r\n")[0])
	if err != nil {
		return &Request{}, fmt.Errorf("error parsing reqeust line err: %s", err)
	}

	return &Request{*rl, nil, nil}, nil

}

func parseRequestLine(ms string) (*RequestLine, error) {
	ln, _, f := strings.Cut(ms, SEPARATOR)
	if f {
		return nil, nil
	}

	parts := strings.Split(ln, " ")

	if len(parts) != 3 {
		return &RequestLine{}, fmt.Errorf("error parsing request line")
	}
	isMatch, err := regexp.MatchString("[A-Z]+", parts[0])
	if err != nil {
		return &RequestLine{}, fmt.Errorf("regex error %w", err)
	}
	if !isMatch {
		return &RequestLine{}, errors.New("invalid method")
	}
	if parts[2] != HTTP_V1 {
		return &RequestLine{}, errors.New("Unsupported HTTP version")
	}

	return &RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   parts[2][5:],
	}, nil

}

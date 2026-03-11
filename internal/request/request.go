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

type parserState string

const (
	Initialized parserState = "initialized"
	Done        parserState = "done"
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
	State       parserState
}

func NewRequest() *Request {
	return &Request{State: Initialized}
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := NewRequest()
	buffer := make([]byte, 1024)
	b_insert_idx := 0
	for request.State != Done {
		n, err := reader.Read(buffer[b_insert_idx:])
		if err != nil {
			return nil, err
		}

		b_insert_idx += n

		_, err = request.parse(buffer[:b_insert_idx+n])
		if err != nil {
			return nil, fmt.Errorf("error parsing the buffer")
		}
		// TODO: implement buffer resizing when max capacity exceeded

	}

	return request, nil

}

func (r *Request) parse(data []byte) (int, error) {
	switch r.State {
	case Initialized:
		rl, n, err := parseRequestLine(string(data))
		if err != nil {
			return 0, err
		}

		if n == 0 {
			return 0, nil
		} else if n > 0 && rl != nil {
			r.RequestLine = *rl
			r.State = Done
		}
		return n, nil

	case Done:
		return 0, fmt.Errorf("Request already done")

	default:
		return 0, fmt.Errorf("Unknown request state")
	}
}

func parseRequestLine(ms string) (*RequestLine, int, error) {
	ln, _, f := strings.Cut(ms, SEPARATOR)
	if !f {
		return nil, 0, nil
	}

	parts := strings.Split(ln, " ")

	if len(parts) != 3 {
		return &RequestLine{}, 0, fmt.Errorf("error parsing request line")
	}
	isMatch, err := regexp.MatchString("[A-Z]+", parts[0])
	if err != nil {
		return &RequestLine{}, 0, fmt.Errorf("regex error %w", err)
	}
	if !isMatch {
		return &RequestLine{}, 0, errors.New("invalid method")
	}
	if parts[2] != HTTP_V1 {
		return &RequestLine{}, 0, errors.New("Unsupported HTTP version")
	}

	return &RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   parts[2][5:],
	}, len(ms), nil

}

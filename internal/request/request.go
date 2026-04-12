package request

import (
	"errors"
	"fmt"
	"http/internal/headers"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"
)

const (
	SEPARATOR = "\r\n"
	HTTP_V1   = "HTTP/1.1"
)

type parserState string

const (
	Initialized parserState = "initialized"
	Headers     parserState = "headers"
	Body        parserState = "body"
	Done        parserState = "done"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte
	State       parserState
}

func NewRequest() *Request {
	return &Request{
		State:   Initialized,
		Headers: *headers.NewHeaders(),
	}
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := NewRequest()
	buffer := make([]byte, 1024)
	b_insert_idx := 0
	b_read_idx := 0
	for request.State != Done {
		n, err := reader.Read(buffer[b_insert_idx:])
		if err != nil {
			if !(err == io.EOF && b_read_idx != b_insert_idx) {
				return nil, err
			}
		}

		b_insert_idx += n

		parsed, err := request.parse(buffer[b_read_idx : b_insert_idx+n])
		b_read_idx += parsed
		if err != nil {
			return nil, fmt.Errorf("error parsing the buffer: %s", err)
		}
		// TODO: implement buffer resizing when maximum capacity exceeded

	}

	return request, nil

}

func (r *Request) parse(data []byte) (int, error) {
	read := 0
	for {
		switch r.State {
		case Initialized:
			rl, n, err := parseRequestLine(string(data[read:]))
			if err != nil {
				return 0, err
			}

			if n == 0 {
				return 0, nil
			} else if n > 0 && rl != nil {
				r.RequestLine = *rl
				r.State = Headers
				read += n

			}

		case Headers:
			n, done, err := r.Headers.Parse(data[read:])
			if err != nil {
				return read, err
			}

			if n == 0 {
				return read, nil
			}
			if done {
				if r.getContentLength() == 0 {
					r.State = Done
				} else {
					r.State = Body
				}
			}
			read += n

		case Body:
			contentlen := r.getContentLength()
			body, n, done := parseRequestBody(data[read:read+contentlen], contentlen)
			if done {
				r.Body = body
				r.State = Done
			}
			return n, nil
		case Done:
			return 0, fmt.Errorf("Request already done")

		default:
			return 0, fmt.Errorf("Unknown request state")
		}
	}
}

func (r *Request) getContentLength() int {

	contentlnStr := r.Headers.Get("Content-Length")
	contentlen, err := strconv.Atoi(contentlnStr)
	if err != nil {
		log.Fatal("Invalid content len")
	}
	return contentlen
}

func parseRequestBody(data []byte, contentLen int) (body []byte, n int, done bool) {
	if len(data) == contentLen {
		return data, contentLen, true
	}
	return []byte{}, 0, false
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
	}, len(ln) + len(SEPARATOR), nil

}

package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaderParsing(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["Host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Valid parse header with already existing headers
	data1 := []byte("Host: localhost:42069\r\n")
	data2 := []byte("Content-Type: text/html\r\n")
	n, done, err = headers.Parse(data1)
	n, done, err = headers.Parse(data2)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "text/html", headers["Content-Type"])
	assert.Equal(t, "localhost:42069", headers["Host"])
	assert.Equal(t, 2, len(headers))

	// Test: Valid done parsing
	headers = NewHeaders()
	data = []byte("\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, 0, len(headers))
	assert.Equal(t, 2, n)
	assert.True(t, done)

	// Test: invalid header spacing
	data = []byte(" Host : localhost:42069\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid header name
	data = []byte("Host@: localhost:42069\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}

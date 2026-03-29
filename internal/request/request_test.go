package request

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type chunkReader struct {
	data            string
	numBytesPerRead int
	pos             int
}

// Read reads up to len(p) or numBytesPerRead bytes from the string per call
// its useful for simulating reading a variable number of bytes per chunk from a network connection
func (cr *chunkReader) Read(p []byte) (n int, err error) {
	if cr.pos >= len(cr.data) {
		return 0, io.EOF
	}
	endIndex := cr.pos + cr.numBytesPerRead
	if endIndex >= len(cr.data) {
		endIndex = len(cr.data)
	}
	n = copy(p, cr.data[cr.pos:endIndex])
	cr.pos += n

	return n, nil
}
func TestGoodRequestLine(t *testing.T) {
	reader := &chunkReader{
		data:            "GET / HTTP/1.1\r\nHost: localhost:42069\r\n\r\n",
		numBytesPerRead: 3,
	}
	r, err := RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)

	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
}

func TestGoodRequestLineWithPath(t *testing.T) {
	reader := &chunkReader{
		data:            "GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\n\r\n",
		numBytesPerRead: 3,
	}
	r, err := RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)

	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
}

func TestGoodPOSTRequestWithPath(t *testing.T) {
	reader := &chunkReader{
		data:            "POST /coffee HTTP/1.1\r\nHost: localhost:42069\r\n\r\n",
		numBytesPerRead: 2,
	}
	r, err := RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)

	assert.Equal(t, "POST", r.RequestLine.Method)
	assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
}

func TestInvalidNumberOfPartsInRequestLine(t *testing.T) {
	tests := []string{
		"GET /coffee\r\nHost: localhost:42069\r\n\r\n",
		"GET /coffee HTTP/1.1 EXTRA\r\nHost: localhost:42069\r\n\r\n",
		"/coffee HTTP/1.1\r\nHost: localhost:42069\r\n\r\n",
	}

	for _, input := range tests {
		_, err := RequestFromReader(strings.NewReader(input))
		require.Error(t, err)
	}
}

func TestInvalidMethodOutOfOrderRequestLine(t *testing.T) {
	tests := []string{
		"get /coffee HTTP/1.1\r\nHost: localhost:42069\r\n\r\n",
		"GeT /coffee HTTP/1.1\r\nHost: localhost:42069\r\n\r\n",
		"G3T /coffee HTTP/1.1\r\nHost: localhost:42069\r\n\r\n",
		"GET! /coffee HTTP/1.1\r\nHost: localhost:42069\r\n\r\n",
	}

	for _, input := range tests {
		_, err := RequestFromReader(strings.NewReader(input))
		require.Error(t, err)
	}
}

func TestInvalidVersionInRequestLine(t *testing.T) {
	tests := []string{
		"GET /coffee HTTP/1.0\r\nHost: localhost:42069\r\n\r\n",
		"GET /coffee HTTP/2.0\r\nHost: localhost:42069\r\n\r\n",
		"GET /coffee HTTPS/1.1\r\nHost: localhost:42069\r\n\r\n",
		"GET /coffee 1.1\r\nHost: localhost:42069\r\n\r\n",
	}

	for _, input := range tests {
		_, err := RequestFromReader(strings.NewReader(input))
		require.Error(t, err)
	}
}

func TestMissingCRLFInRequestLine(t *testing.T) {
	_, err := RequestFromReader(strings.NewReader(
		"GET /coffee HTTP/1.1",
	))
	require.Error(t, err)
}

func TestEmptyInput(t *testing.T) {
	_, err := RequestFromReader(strings.NewReader(""))
	require.Error(t, err)
}

func TestIgnoresEverythingAfterRequestLine(t *testing.T) {
	r, err := RequestFromReader(strings.NewReader(
		"GET /tea HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\nhello",
	))
	require.NoError(t, err)
	require.NotNil(t, r)

	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/tea", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
}

package request

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGoodRequestLine(t *testing.T) {
	r, err := RequestFromReader(strings.NewReader(
		"GET / HTTP/1.1\r\nHost: localhost:42069\r\n\r\n",
	))
	require.NoError(t, err)
	require.NotNil(t, r)

	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
}

func TestGoodRequestLineWithPath(t *testing.T) {
	r, err := RequestFromReader(strings.NewReader(
		"GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\n\r\n",
	))
	require.NoError(t, err)
	require.NotNil(t, r)

	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
}

func TestGoodPOSTRequestWithPath(t *testing.T) {
	r, err := RequestFromReader(strings.NewReader(
		"POST /coffee HTTP/1.1\r\nHost: localhost:42069\r\n\r\n",
	))
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

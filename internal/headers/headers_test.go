package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidSingleHeader(t *testing.T) {
	headers := Headers{}
	data := []byte("HoSt: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)
}

func TestValidSingleHeaderWithExtraSpace(t *testing.T) {
	headers := Headers{}
	data := []byte(" Host: localhost:42069 \r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 25, n)
	assert.False(t, done)
}

func TestValid2HeadersWithExistingHeaders(t *testing.T) {
	headers := Headers{}

	data := []byte("Host: localhost\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, 17, n)
	assert.False(t, done)
	assert.Equal(t, "localhost", headers["host"])

	data = []byte("User-Agent: test\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, 18, n)
	assert.False(t, done)
	assert.Equal(t, "test", headers["user-agent"])

	assert.Equal(t, "localhost", headers["host"])
}

func TestValidHeaderMultipleValues(t *testing.T) {
	headers := Headers{"host": "first.com"}
	data := []byte("HoSt: second.com\r\n\r\n")
	_, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "first.com,second.com", headers["host"])
	assert.False(t, done)
}

func TestInvalidSpacingHeader(t *testing.T) {
	headers := Headers{}
	data := []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err := headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}

func TestInvalidCharacterHeader(t *testing.T) {
	headers := Headers{}
	data := []byte("H@st: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.Error(t, err)
	require.Empty(t, headers)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}

func TestValidDone(t *testing.T) {
	headers := Headers{}
	data := []byte("\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.Empty(t, headers)
	assert.Equal(t, 2, n)
	assert.True(t, done)
}

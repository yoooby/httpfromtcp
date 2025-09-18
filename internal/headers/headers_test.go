package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)




func TestHeaders(t *testing.T) {
		// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid Character in fieldname
	headers = NewHeaders()
	data = []byte("       H©st : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// 	// Test: Emtpy CRLF
	headers = NewHeaders()
	data = []byte("\r\n")
	n, done, err = headers.Parse(data)
	require.Equal(t, n, 2)
	require.NoError(t, err)
	require.True(t, done)

	// Test: field name with multiple values
    headers = NewHeaders()
    data = []byte("Host: localhost:42069\r\n\r\n")
    n, done, err = headers.Parse(data)
    require.NoError(t, err)
    require.False(t, done)
    assert.Equal(t, 23, n)
    assert.Equal(t, "localhost:42069", headers["host"])
    data = []byte("Host: zebi:2016\r\n\r\n")
    n, done, err = headers.Parse(data)
    require.NoError(t, err)
    require.False(t, done)
    assert.Equal(t, 17, n)
    assert.Equal(t, "localhost:42069, zebi:2016", headers["host"])

}


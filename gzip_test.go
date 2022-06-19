package compresshandler

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/klauspost/compress/gzip"
	"github.com/stretchr/testify/assert"
)

func gzipSlice(a []byte) ([]byte, error) {
	var b bytes.Buffer
	gzipWriter := gzip.NewWriter(&b)
	if _, err := gzipWriter.Write(a); err != nil {
		panic("gzipSlice Write error: " + err.Error())
	}

	if err := gzipWriter.Flush(); err != nil {
		panic("gzipSlice Flush error: " + err.Error())
	}

	return b.Bytes(), nil
}

func ungzipSlice(a []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewBuffer(a))
	if err != nil {
		panic("ungzipSlice " + err.Error())
	}
	retval, _ := io.ReadAll(r)

	return retval, nil
}

func TestGzip(t *testing.T) {
	test := "there is fake string *^(^$&*^&"
	input := []byte(test)
	gzipped, err := gzipSlice(input)
	if err != nil {
		panic("gzipSlice " + err.Error())
	}

	ungzipped, err := ungzipSlice(gzipped)
	if err != nil {
		panic("ungzipSlice: " + err.Error())
	}

	assert.True(t, strings.TrimSpace(string(input)) == strings.TrimSpace(string(ungzipped)))
}

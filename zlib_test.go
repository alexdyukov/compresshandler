package compresshandler

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/klauspost/compress/zlib"
	"github.com/stretchr/testify/assert"
)

func zlibSlice(a []byte) ([]byte, error) {
	var b bytes.Buffer
	zlibWriter := zlib.NewWriter(&b)
	if _, err := zlibWriter.Write(a); err != nil {
		panic("gzipSlice Write error: " + err.Error())
	}

	if err := zlibWriter.Flush(); err != nil {
		panic("gzipSlice Flush error: " + err.Error())
	}

	return b.Bytes(), nil
}

func unzlibSlice(a []byte) ([]byte, error) {
	r, err := zlib.NewReader(bytes.NewBuffer(a))
	if err != nil {
		panic("unzlibSlice " + err.Error())
	}
	retval, _ := io.ReadAll(r)

	return retval, nil
}

func TestZlib(t *testing.T) {
	test := "there is fake string *^(^$&*^&"
	input := []byte(test)
	zlibbed, err := zlibSlice(input)
	if err != nil {
		panic("zlibSlice " + err.Error())
	}

	unzlibbed, err := unzlibSlice(zlibbed)
	if err != nil {
		panic("unzlibSlice: " + err.Error())
	}

	assert.True(t, strings.TrimSpace(string(input)) == strings.TrimSpace(string(unzlibbed)))
}

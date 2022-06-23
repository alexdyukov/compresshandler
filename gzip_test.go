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
		return b.Bytes(), err
	}

	if err := gzipWriter.Flush(); err != nil {
		return b.Bytes(), err
	}

	return b.Bytes(), nil
}

func ungzipSlice(a []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewBuffer(a))
	if err != nil {
		return nil, err
	}
	retval, _ := io.ReadAll(r)

	return retval, nil
}

func TestGzip(t *testing.T) {
	test := "there is fake string *^(^$&*^&"
	input := []byte(test)
	gziped, err := gzipSlice(input)
	if err != nil {
		t.Fatalf("gzipSlice: %v", err)
	}

	ungziped, err := ungzipSlice(gziped)
	if err != nil {
		t.Fatalf("ungzipSlice: %v", err)
	}

	assert.True(t, strings.TrimSpace(string(input)) == strings.TrimSpace(string(ungziped)))
}

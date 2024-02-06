package compressor_test

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/alexdyukov/compresshandler/internal/compressor"
	"github.com/stretchr/testify/assert"
)

func stdGzipDecompress(to *bytes.Buffer, from io.Reader) error {
	reader, err := gzip.NewReader(from)
	if err != nil {
		return fmt.Errorf("compressor: gzip_test: failed to initialize reader: %w", err)
	}

	_, err = to.ReadFrom(reader)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		return fmt.Errorf("compressor: gzip_test: failed to read data: %w", err)
	}

	return nil
}

func TestGzip(t *testing.T) {
	t.Parallel()

	testInput := []byte("there is fake string *^(^$&*^&")

	gziped := &bytes.Buffer{}
	ungziped := &bytes.Buffer{}

	compressor := compressor.NewGzip((gzip.BestCompression - gzip.BestSpeed) / 2)

	if err := compressor.Compress(gziped, testInput); err != nil {
		t.Fatalf("TestGzipCompression: Compress: %v of type %T", err, err)
	}

	if err := stdGzipDecompress(ungziped, gziped); err != nil {
		t.Fatalf("TestGzipCompression: stdGzipDecompress: %v of type %T", err, err)
	}

	assert.True(t, strings.TrimSpace(string(testInput)) == strings.TrimSpace(ungziped.String()))
}

func BenchmarkGzipCompressor(b *testing.B) {
	testInput := []byte("there is fake string *^(^$&*^&")
	gziped := &bytes.Buffer{}
	compressor := compressor.NewGzip((gzip.BestCompression - gzip.BestSpeed) / 2)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := compressor.Compress(gziped, testInput); err != nil {
			b.FailNow()
		}

		gziped.Reset()
	}
}

package decompressor_test

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/alexdyukov/compresshandler/internal/decompressor"
	"github.com/stretchr/testify/assert"
)

func stdGzipCompress(level int, to io.Writer, from *bytes.Buffer) error {
	writer, err := gzip.NewWriterLevel(to, level)
	if err != nil {
		return fmt.Errorf("decompressor: gzip_test: failed to initialize writer: %w", err)
	}

	if _, err := writer.Write(from.Bytes()); err != nil {
		return fmt.Errorf("decompressor: gzip_test: failed to write data: %w", err)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("decompressor: gzip_test: failed to flush data: %w", err)
	}

	return nil
}

func TestGzip(t *testing.T) {
	t.Parallel()

	testInput := bytes.NewBufferString("there is fake string *^(^$&*^&")

	gziped := &bytes.Buffer{}
	ungziped := &bytes.Buffer{}

	decompressor := decompressor.NewGzip()

	if err := stdGzipCompress((gzip.BestCompression-gzip.BestSpeed)/2, gziped, testInput); err != nil {
		t.Fatalf("TestGzipDecompression: stdGzipCompress: %v of type %T", err, err)
	}

	if err := decompressor.Decompress(ungziped, gziped); err != nil {
		t.Fatalf("TestGzipDecompression: Decompress: %v of type %T", err, err)
	}

	assert.True(t, strings.TrimSpace(testInput.String()) == strings.TrimSpace(ungziped.String()))
}

func BenchmarkGzip(b *testing.B) {
	testInput := bytes.NewBufferString("there is fake string *^(^$&*^&")
	gziped := &bytes.Buffer{}
	ungziped := &bytes.Buffer{}
	decompressor := decompressor.NewGzip()
	level := (gzip.BestCompression - gzip.BestSpeed) / 2

	if err := stdGzipCompress(level, gziped, testInput); err != nil {
		b.Fatalf("BenchmarkGzipDedecompressor: Compress: %v of type %T", err, err)
	}

	data := gziped.Bytes()
	reader := bytes.NewReader(data)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := decompressor.Decompress(ungziped, reader); err != nil {
			b.FailNow()
		}

		reader.Reset(data)
		ungziped.Reset()
	}
}

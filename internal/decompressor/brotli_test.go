package decompressor_test

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/alexdyukov/compresshandler/internal/decompressor"
	"github.com/andybalholm/brotli"
	"github.com/stretchr/testify/assert"
)

func stdBrotliCompress(level int, target io.Writer, from *bytes.Buffer) error {
	writer := brotli.NewWriterLevel(target, level)

	if _, err := writer.Write(from.Bytes()); err != nil {
		return fmt.Errorf("decompressor: brotli_test: failed to write data: %w", err)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("decompressor: brotli_test: failed to flush data: %w", err)
	}

	return nil
}

func TestBrotli(t *testing.T) {
	t.Parallel()

	testInput := bytes.NewBufferString("there is fake string *^(^$&*^&")

	brotlied := &bytes.Buffer{}
	unbrotlied := &bytes.Buffer{}

	decompressor := decompressor.NewBrotli()

	if err := stdBrotliCompress((brotli.BestCompression-brotli.BestSpeed)/2, brotlied, testInput); err != nil {
		t.Fatalf("TestBrotliDecompression: stdBrotliCompress: %v of type %T", err, err)
	}

	if err := decompressor.Decompress(unbrotlied, brotlied); err != nil {
		t.Fatalf("TestBrotliDecompression: Decompress: %v of type %T", err, err)
	}

	assert.True(t, strings.TrimSpace(testInput.String()) == strings.TrimSpace(unbrotlied.String()))
}

func BenchmarkBrotli(b *testing.B) {
	testInput := bytes.NewBufferString("there is fake string *^(^$&*^&")
	brotlied := &bytes.Buffer{}
	unbrotlied := &bytes.Buffer{}
	decompressor := decompressor.NewBrotli()
	level := (brotli.BestCompression - brotli.BestSpeed) / 2

	if err := stdBrotliCompress(level, brotlied, testInput); err != nil {
		b.Fatalf("BenchmarkBrotliDedecompressor: Compress: %v of type %T", err, err)
	}

	data := brotlied.Bytes()
	reader := bytes.NewReader(data)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := decompressor.Decompress(unbrotlied, reader); err != nil {
			b.FailNow()
		}

		reader.Reset(data)
		unbrotlied.Reset()
	}
}

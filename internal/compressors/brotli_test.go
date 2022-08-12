package compressors_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/alexdyukov/compresshandler/internal/compressors"
	"github.com/andybalholm/brotli"
	"github.com/stretchr/testify/assert"
)

func stdBrotliDecompress(target *bytes.Buffer, from io.Reader) error {
	reader := brotli.NewReader(from)

	_, err := target.ReadFrom(reader)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		return fmt.Errorf("compressors: brotli_test: failed to read data: %w", err)
	}

	return nil
}

func TestBrotli(t *testing.T) {
	t.Parallel()

	testInput := bytes.NewBufferString("there is fake string *^(^$&*^&")

	brotlied := &bytes.Buffer{}
	unbrotlied := &bytes.Buffer{}

	compressor := compressors.NewBrotli((brotli.BestCompression - brotli.BestSpeed) / 2)

	if err := compressor.Compress(brotlied, testInput); err != nil {
		t.Fatalf("TestBrotliCompression: Compress: %v of type %T", err, err)
	}

	if err := stdBrotliDecompress(unbrotlied, brotlied); err != nil {
		t.Fatalf("TestBrotliCompression: stdBrotliDecompress: %v of type %T", err, err)
	}

	assert.True(t, strings.TrimSpace(testInput.String()) == strings.TrimSpace(unbrotlied.String()))
}

func BenchmarkBrotli(b *testing.B) {
	testInput := bytes.NewBufferString("there is fake string *^(^$&*^&")
	brotlied := &bytes.Buffer{}
	compressor := compressors.NewBrotli((brotli.BestCompression - brotli.BestSpeed) / 2)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := compressor.Compress(brotlied, testInput); err != nil {
			b.FailNow()
		}

		brotlied.Reset()
	}
}

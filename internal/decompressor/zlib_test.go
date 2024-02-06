package decompressor_test

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/alexdyukov/compresshandler/internal/decompressor"
	"github.com/stretchr/testify/assert"
)

func stdZlibCompress(level int, to io.Writer, from *bytes.Buffer) error {
	writer, err := zlib.NewWriterLevel(to, level)
	if err != nil {
		return fmt.Errorf("decompressor: zlib_test: failed to initialize writer: %w", err)
	}

	if _, err := writer.Write(from.Bytes()); err != nil {
		return fmt.Errorf("decompressor: zlib_test: failed to write data: %w", err)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("decompressor: zlib_test: failed to flush data: %w", err)
	}

	return nil
}

func TestZlib(t *testing.T) {
	t.Parallel()

	testInput := bytes.NewBufferString("there is fake string *^(^$&*^&")

	zlibed := &bytes.Buffer{}
	unzlibed := &bytes.Buffer{}

	decompressor := decompressor.NewZlib()

	if err := stdZlibCompress((zlib.BestCompression-zlib.BestSpeed)/2, zlibed, testInput); err != nil {
		t.Fatalf("TestZlibDecompression: stdZlibCompress: %v of type %T", err, err)
	}

	if err := decompressor.Decompress(unzlibed, zlibed); err != nil {
		t.Fatalf("TestZlibDecompression: Decompress: %v of type %T", err, err)
	}

	assert.True(t, strings.TrimSpace(testInput.String()) == strings.TrimSpace(unzlibed.String()))
}

func BenchmarkZlib(b *testing.B) {
	testInput := bytes.NewBufferString("there is fake string *^(^$&*^&")
	zlibed := &bytes.Buffer{}
	unzlibed := &bytes.Buffer{}
	decompressor := decompressor.NewZlib()
	level := (zlib.BestCompression - zlib.BestSpeed) / 2

	if err := stdZlibCompress(level, zlibed, testInput); err != nil {
		b.Fatalf("BenchmarkZlibDedecompressor: Compress: %v of type %T", err, err)
	}

	data := zlibed.Bytes()
	reader := bytes.NewReader(data)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := decompressor.Decompress(unzlibed, reader); err != nil {
			b.FailNow()
		}

		reader.Reset(data)
		unzlibed.Reset()
	}
}

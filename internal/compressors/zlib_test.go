package compressors_test

import (
	"bytes"
	"compress/zlib"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/alexdyukov/compresshandler/internal/compressors"
	"github.com/stretchr/testify/assert"
)

func stdZlibDecompress(target *bytes.Buffer, from io.Reader) error {
	reader, err := zlib.NewReader(from)
	if err != nil {
		return fmt.Errorf("compressors: zlib_test: failed to initialize reader: %w", err)
	}

	_, err = target.ReadFrom(reader)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		return fmt.Errorf("compressors: zlib_test: failed to read data: %w", err)
	}

	return nil
}

func TestZlib(t *testing.T) {
	t.Parallel()

	testInput := bytes.NewBufferString("there is fake string *^(^$&*^&")

	zlibed := &bytes.Buffer{}
	unzlibed := &bytes.Buffer{}

	compressor := compressors.NewZlib((zlib.BestCompression - zlib.BestSpeed) / 2)

	if err := compressor.Compress(zlibed, testInput); err != nil {
		t.Fatalf("TestZlibCompression: Compress: %v of type %T", err, err)
	}

	if err := stdZlibDecompress(unzlibed, zlibed); err != nil {
		t.Fatalf("TestZlibCompression: stdZlibDecompress: %v of type %T", err, err)
	}

	assert.True(t, strings.TrimSpace(testInput.String()) == strings.TrimSpace(unzlibed.String()))
}

func BenchmarkZlib(b *testing.B) {
	testInput := bytes.NewBufferString("there is fake string *^(^$&*^&")
	zlibed := &bytes.Buffer{}
	compressor := compressors.NewZlib((zlib.BestCompression - zlib.BestSpeed) / 2)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := compressor.Compress(zlibed, testInput); err != nil {
			b.FailNow()
		}

		zlibed.Reset()
	}
}

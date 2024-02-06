package decompressor

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/andybalholm/brotli"
)

// Brotli is brotli typed (http encoding type "br") Decompressor.
type Brotli struct {
	pool sync.Pool
}

func getInitialBrotliSlice() []byte {
	initialBuffer := bytes.Buffer{}

	writer := brotli.NewWriterLevel(&initialBuffer, brotli.BestSpeed)

	if err := writer.Close(); err != nil {
		panic("unreachable code")
	}

	return initialBuffer.Bytes()
}

// NewBrotli creates new brotli typed (http encoding type "br") Decompressor.
func NewBrotli() *Brotli {
	initialSlice := getInitialBrotliSlice()

	return &Brotli{sync.Pool{
		New: func() interface{} {
			return brotli.NewReader(bytes.NewBuffer(initialSlice))
		},
	}}
}

// Decompress decompressing bytes from src to dst with brotli decompress algo until error occurs or end of src.
func (decompressor *Brotli) Decompress(dst *bytes.Buffer, src io.Reader) error {
	reader, ok := decompressor.pool.Get().(*brotli.Reader)
	if !ok {
		panic("unreachable code")
	}

	defer decompressor.pool.Put(reader)

	if err := reader.Reset(src); err != nil {
		return fmt.Errorf("decompressor: brotli: failed to initialize reader from pool: %w", err)
	}

	_, err := dst.ReadFrom(reader)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		return fmt.Errorf("decompressor: brotli: failed to decompress data: %w", err)
	}

	return nil
}

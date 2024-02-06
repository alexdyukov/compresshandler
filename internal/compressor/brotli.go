package compressor

import (
	"fmt"
	"io"
	"sync"

	"github.com/andybalholm/brotli"
)

// Brotli is brotli typed (http encoding type "br") Compressor.
type Brotli struct {
	pool sync.Pool
}

const (
	BrotliBestCompression = brotli.BestCompression
	BrotliBestSpeed       = brotli.BestSpeed
)

// NewBrotli creates new brotli typed (http encoding type "br") Compressor.
func NewBrotli(level int) *Brotli {
	return &Brotli{sync.Pool{
		New: func() interface{} {
			return brotli.NewWriterLevel(nil, level)
		},
	}}
}

// Compress compressing bytes from src to dst with brotli compressing algo until error occurs or end of src.
func (compressor *Brotli) Compress(dst io.Writer, src []byte) error {
	writer, ok := compressor.pool.Get().(*brotli.Writer)
	if !ok {
		panic("unreachable code")
	}

	defer compressor.pool.Put(writer)

	writer.Reset(dst)

	if _, err := writer.Write(src); err != nil {
		return fmt.Errorf("compressor: brotli: failed to write data: %w", err)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("compressor: brotli: failed to flush data: %w", err)
	}

	return nil
}

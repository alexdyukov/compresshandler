package compressor

import (
	"fmt"
	"io"
	"sync"

	"github.com/klauspost/compress/gzip"
)

// Gzip is gzip typed (http encoding type "gzip") Compressor.
type Gzip struct {
	pool sync.Pool
}

const (
	GzipBestCompression = gzip.BestCompression
	GzipBestSpeed       = gzip.BestSpeed
)

// NewGzip creates new gzip typed (http encoding type "gzip") Compressor.
func NewGzip(level int) *Gzip {
	return &Gzip{sync.Pool{
		New: func() interface{} {
			writer, err := gzip.NewWriterLevel(nil, level)
			if err != nil {
				panic("unreachable code")
			}

			return writer
		},
	}}
}

// Compress compressing bytes from src to dst with gzip compressing algo until error occurs or end of src.
func (compressorPool *Gzip) Compress(dst io.Writer, src []byte) error {
	writer, ok := compressorPool.pool.Get().(*gzip.Writer)
	if !ok {
		panic("unreachable code")
	}

	defer compressorPool.pool.Put(writer)

	writer.Reset(dst)

	if _, err := writer.Write(src); err != nil {
		return fmt.Errorf("compressor: gzip: failed to write data: %w", err)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("compressor: gzip: failed to flush data: %w", err)
	}

	return nil
}

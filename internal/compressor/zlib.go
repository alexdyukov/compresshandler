package compressor

import (
	"fmt"
	"io"
	"sync"

	"github.com/klauspost/compress/zlib"
)

// Zlib is zlib typed (http encoding type "deflate") Compressor.
type Zlib struct {
	pool sync.Pool
}

const (
	ZlibBestCompression = zlib.BestCompression
	ZlibBestSpeed       = zlib.BestSpeed
)

// NewZlib creates new zlib typed (http encoding type "deflate") Compressor.
func NewZlib(level int) *Zlib {
	return &Zlib{sync.Pool{
		New: func() interface{} {
			writer, err := zlib.NewWriterLevel(nil, level)
			if err != nil {
				panic("unreachable code")
			}

			return writer
		},
	}}
}

// Compress compressing bytes from src to dst with zlib compressing algo until error occurs or end of src.
func (compressorPool *Zlib) Compress(dst io.Writer, src []byte) error {
	writer, ok := compressorPool.pool.Get().(*zlib.Writer)
	if !ok {
		panic("unreachable code")
	}

	defer compressorPool.pool.Put(writer)

	writer.Reset(dst)

	if _, err := writer.Write(src); err != nil {
		return fmt.Errorf("compressor: zlib: failed to write data: %w", err)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("compressor: zlib: failed to flush data: %w", err)
	}

	return nil
}

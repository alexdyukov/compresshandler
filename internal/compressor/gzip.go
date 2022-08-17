package compressor

import (
	"bytes"
	"fmt"
	"io"
	"sync"

	"github.com/klauspost/compress/gzip"
)

type Gzip struct {
	sync.Pool
}

const (
	GzipBestCompression = gzip.BestCompression
	GzipBestSpeed       = gzip.BestSpeed
)

func NewGzip(level int) *Gzip {
	return &Gzip{sync.Pool{
		New: func() interface{} {
			writer, err := gzip.NewWriterLevel(&bytes.Buffer{}, level)
			if err != nil {
				panic("unreachable code")
			}

			return writer
		},
	}}
}

func (compressorPool *Gzip) Compress(target io.Writer, from []byte) error {
	writer, ok := compressorPool.Get().(*gzip.Writer)
	if !ok {
		panic("unreachable code")
	}

	defer compressorPool.Put(writer)

	writer.Reset(target)

	if _, err := writer.Write(from); err != nil {
		return fmt.Errorf("compressor: gzip: failed to write data: %w", err)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("compressor: gzip: failed to flush data: %w", err)
	}

	return nil
}

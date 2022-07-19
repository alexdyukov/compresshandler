package compressors

import (
	"bytes"
	"fmt"
	"io"
	"sync"

	"github.com/klauspost/compress/gzip"
)

type Gzip [gzip.BestCompression - gzip.BestSpeed]*sync.Pool

func NewGzip() *Gzip {
	var compressor Gzip

	for i := range compressor {
		i := i
		compressor[i] = &sync.Pool{
			New: func() interface{} {
				writer, err := gzip.NewWriterLevel(&bytes.Buffer{}, i)
				if err != nil {
					panic("unreachable code")
				}

				return writer
			},
		}
	}

	return &compressor
}

func (compressor *Gzip) Compress(level int, target io.Writer, from *bytes.Buffer) error {
	pool := compressor[level-gzip.BestSpeed]

	writer, ok := pool.Get().(*gzip.Writer)
	if !ok {
		panic("unreachable code")
	}

	defer pool.Put(writer)

	writer.Reset(target)

	if _, err := writer.Write(from.Bytes()); err != nil {
		return fmt.Errorf("compressors: gzip: failed to write data: %w", err)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("compressors: gzip: failed to flush data: %w", err)
	}

	return nil
}

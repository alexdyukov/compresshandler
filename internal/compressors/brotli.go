package compressors

import (
	"bytes"
	"fmt"
	"io"
	"sync"

	"github.com/andybalholm/brotli"
)

type Brotli [brotli.BestCompression - brotli.BestSpeed]*sync.Pool

func NewBrotli() *Brotli {
	var compressor Brotli

	for i := range compressor {
		i := i
		compressor[i] = &sync.Pool{
			New: func() interface{} {
				return brotli.NewWriterLevel(&bytes.Buffer{}, i)
			},
		}
	}

	return &compressor
}

func (compressor *Brotli) Compress(level int, target io.Writer, from *bytes.Buffer) error {
	pool := compressor[level-brotli.BestSpeed]

	writer, ok := pool.Get().(*brotli.Writer)
	if !ok {
		panic("unreachable code")
	}

	defer pool.Put(writer)

	writer.Reset(target)

	if _, err := writer.Write(from.Bytes()); err != nil {
		return fmt.Errorf("compressors: brotli: failed to write data: %w", err)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("compressors: brotli: failed to flush data: %w", err)
	}

	return nil
}

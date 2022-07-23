package compressors

import (
	"bytes"
	"fmt"
	"io"
	"sync"

	"github.com/klauspost/compress/zlib"
)

type Zlib [zlib.BestCompression - zlib.BestSpeed]*sync.Pool

func NewZlib() *Zlib {
	var compressor Zlib

	for i := range compressor {
		i := i
		compressor[i] = &sync.Pool{
			New: func() interface{} {
				writer, err := zlib.NewWriterLevel(&bytes.Buffer{}, i)
				if err != nil {
					panic("unreachable code")
				}

				return writer
			},
		}
	}

	return &compressor
}

func (compressor *Zlib) Compress(level int, target io.Writer, from *bytes.Buffer) error {
	pool := compressor[level-zlib.BestSpeed]

	writer, ok := pool.Get().(*zlib.Writer)
	if !ok {
		panic("unreachable code")
	}

	defer pool.Put(writer)

	writer.Reset(target)

	if _, err := writer.Write(from.Bytes()); err != nil {
		return fmt.Errorf("compressors: zlib: failed to write data: %w", err)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("compressors: zlib: failed to flush data: %w", err)
	}

	return nil
}

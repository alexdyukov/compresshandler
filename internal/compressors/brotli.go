package compressors

import (
	"bytes"
	"fmt"
	"io"
	"sync"

	"github.com/andybalholm/brotli"
)

type Brotli struct {
	sync.Pool
}

func NewBrotli(level int) *Brotli {
	return &Brotli{sync.Pool{
		New: func() interface{} {
			return brotli.NewWriterLevel(&bytes.Buffer{}, level)
		},
	}}
}

func (compressorPool *Brotli) Compress(target io.Writer, from *bytes.Buffer) error {
	writer, ok := compressorPool.Get().(*brotli.Writer)
	if !ok {
		panic("unreachable code")
	}

	defer compressorPool.Put(writer)

	writer.Reset(target)

	if _, err := writer.Write(from.Bytes()); err != nil {
		return fmt.Errorf("compressors: brotli: failed to write data: %w", err)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("compressors: brotli: failed to flush data: %w", err)
	}

	return nil
}

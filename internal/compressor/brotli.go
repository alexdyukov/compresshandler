package compressor

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

const (
	BrotliBestCompression = brotli.BestCompression
	BrotliBestSpeed       = brotli.BestSpeed
)

func NewBrotli(level int) *Brotli {
	return &Brotli{sync.Pool{
		New: func() interface{} {
			return brotli.NewWriterLevel(&bytes.Buffer{}, level)
		},
	}}
}

func (compressorPool *Brotli) Compress(target io.Writer, from []byte) error {
	writer, ok := compressorPool.Get().(*brotli.Writer)
	if !ok {
		panic("unreachable code")
	}

	defer compressorPool.Put(writer)

	writer.Reset(target)

	if _, err := writer.Write(from); err != nil {
		return fmt.Errorf("compressor: brotli: failed to write data: %w", err)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("compressor: brotli: failed to flush data: %w", err)
	}

	return nil
}

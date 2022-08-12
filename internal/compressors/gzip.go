package compressors

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

func (compressorPool *Gzip) Compress(target io.Writer, from *bytes.Buffer) error {
	writer, ok := compressorPool.Get().(*gzip.Writer)
	if !ok {
		panic("unreachable code")
	}

	defer compressorPool.Put(writer)

	writer.Reset(target)

	if _, err := writer.Write(from.Bytes()); err != nil {
		return fmt.Errorf("compressors: gzip: failed to write data: %w", err)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("compressors: gzip: failed to flush data: %w", err)
	}

	return nil
}

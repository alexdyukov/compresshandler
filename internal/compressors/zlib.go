package compressors

import (
	"bytes"
	"fmt"
	"io"
	"sync"

	"github.com/klauspost/compress/zlib"
)

type Zlib struct {
	sync.Pool
}

func NewZlib(level int) *Zlib {
	return &Zlib{sync.Pool{
		New: func() interface{} {
			writer, err := zlib.NewWriterLevel(&bytes.Buffer{}, level)
			if err != nil {
				panic("unreachable code")
			}

			return writer
		},
	}}
}

func (compressorPool *Zlib) Compress(target io.Writer, from *bytes.Buffer) error {
	writer, ok := compressorPool.Get().(*zlib.Writer)
	if !ok {
		panic("unreachable code")
	}

	defer compressorPool.Put(writer)

	writer.Reset(target)

	if _, err := writer.Write(from.Bytes()); err != nil {
		return fmt.Errorf("compressors: zlib: failed to write data: %w", err)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("compressors: zlib: failed to flush data: %w", err)
	}

	return nil
}

package decompressor

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/klauspost/compress/gzip"
)

type Gzip struct {
	sync.Pool
}

func getInitialGzipSlice() []byte {
	initialBuffer := bytes.Buffer{}

	writer, err := gzip.NewWriterLevel(&initialBuffer, gzip.BestSpeed)
	if err != nil {
		panic("unreachable code")
	}

	if err := writer.Close(); err != nil {
		panic("unreachable code")
	}

	return initialBuffer.Bytes()
}

func NewGzip() *Gzip {
	initialSlice := getInitialGzipSlice()

	return &Gzip{sync.Pool{
		New: func() interface{} {
			reader, err := gzip.NewReader(bytes.NewBuffer(initialSlice))
			if err != nil {
				panic("unreachable code")
			}

			return reader
		},
	}}
}

func (decompressor *Gzip) Decompress(target *bytes.Buffer, from io.Reader) error {
	reader, ok := decompressor.Get().(*gzip.Reader)
	if !ok {
		panic("unreachable code")
	}

	defer decompressor.Put(reader)

	if err := reader.Reset(from); err != nil {
		return fmt.Errorf("decompressor: gzip: failed to initialize reader from pool: %w", err)
	}

	_, err := target.ReadFrom(reader)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		return fmt.Errorf("decompressor: gzip: failed to decompress data: %w", err)
	}

	return nil
}

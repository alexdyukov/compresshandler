package decompressors

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/andybalholm/brotli"
)

type Brotli struct {
	sync.Pool
}

func getInitialBrotliSlice() []byte {
	initialBuffer := bytes.Buffer{}

	writer := brotli.NewWriterLevel(&initialBuffer, brotli.BestSpeed)

	if err := writer.Close(); err != nil {
		panic("unreachable code")
	}

	return initialBuffer.Bytes()
}

func NewBrotli() *Brotli {
	initialSlice := getInitialBrotliSlice()

	return &Brotli{sync.Pool{
		New: func() interface{} {
			return brotli.NewReader(bytes.NewBuffer(initialSlice))
		},
	}}
}

func (decompressor *Brotli) Decompress(target *bytes.Buffer, from io.Reader) error {
	reader, ok := decompressor.Get().(*brotli.Reader)
	if !ok {
		panic("unreachable code")
	}

	defer decompressor.Put(reader)

	if err := reader.Reset(from); err != nil {
		return fmt.Errorf("decompressors: brotli: failed to initialize reader from pool: %w", err)
	}

	_, err := target.ReadFrom(reader)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		return fmt.Errorf("decompressors: brotli: failed to decompress data: %w", err)
	}

	return nil
}

package decompressor

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/klauspost/compress/gzip"
)

// Gzip is gzip typed (http encoding type "gzip") Decompressor.
type Gzip struct {
	pool sync.Pool
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

// NewGzip creates new gzip typed (http encoding type "gzip") Decompressor.
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

// Decompress decompressing bytes from src to dst with gzip decompress algo until error occurs or end of src.
func (decompressor *Gzip) Decompress(dst *bytes.Buffer, src io.Reader) error {
	reader, ok := decompressor.pool.Get().(*gzip.Reader)
	if !ok {
		panic("unreachable code")
	}

	defer decompressor.pool.Put(reader)

	if err := reader.Reset(src); err != nil {
		return fmt.Errorf("decompressor: gzip: failed to initialize reader from pool: %w", err)
	}

	_, err := dst.ReadFrom(reader)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		return fmt.Errorf("decompressor: gzip: failed to decompress data: %w", err)
	}

	return nil
}

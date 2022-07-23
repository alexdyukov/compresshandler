package decompressors

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/klauspost/compress/zlib"
)

type (
	zlibReader interface {
		io.ReadCloser
		Reset(r io.Reader, dict []byte) error
	}
	Zlib struct {
		sync.Pool
	}
)

func getInitialZlibSlice() []byte {
	initialBuffer := bytes.Buffer{}

	writer, err := zlib.NewWriterLevel(&initialBuffer, zlib.BestSpeed)
	if err != nil {
		panic("unreachable code")
	}

	if err := writer.Close(); err != nil {
		panic("unreachable code")
	}

	return initialBuffer.Bytes()
}

func NewZlib() *Zlib {
	initialSlice := getInitialZlibSlice()

	return &Zlib{sync.Pool{
		New: func() interface{} {
			zReader, err := zlib.NewReader(bytes.NewBuffer(initialSlice))
			if err != nil {
				panic("unreachable code")
			}

			reader, ok := zReader.(zlibReader)
			if !ok {
				panic("unreachable code")
			}

			return reader
		},
	}}
}

func (decompressor *Zlib) Decompress(target *bytes.Buffer, from io.Reader) error {
	reader, ok := decompressor.Get().(zlibReader)
	if !ok {
		panic("unreachable code")
	}

	defer decompressor.Put(reader)

	if err := reader.Reset(from, nil); err != nil {
		return fmt.Errorf("decompressors: zlib: failed to initialize reader from pool: %w", err)
	}

	_, err := target.ReadFrom(reader)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		return fmt.Errorf("decompressors: zlib: failed to decompress data: %w", err)
	}

	return nil
}

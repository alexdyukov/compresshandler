package decompressor

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
	// Zlib is zlib typed (http encoding type "deflate") Decompressor.
	Zlib struct {
		pool sync.Pool
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

// NewZlib creates new zlib typed (http encoding type "deflate") Decompressor.
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

// Decompress decompressing bytes from src to dst with zlib decompress algo until error occurs or end of src.
func (decompressor *Zlib) Decompress(dst *bytes.Buffer, src io.Reader) error {
	reader, ok := decompressor.pool.Get().(zlibReader)
	if !ok {
		panic("unreachable code")
	}

	defer decompressor.pool.Put(reader)

	if err := reader.Reset(src, nil); err != nil {
		return fmt.Errorf("decompressor: zlib: failed to initialize reader from pool: %w", err)
	}

	_, err := dst.ReadFrom(reader)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		return fmt.Errorf("decompressor: zlib: failed to decompress data: %w", err)
	}

	return nil
}

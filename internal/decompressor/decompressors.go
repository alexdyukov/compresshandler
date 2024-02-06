package decompressor

import (
	"bytes"
	"io"
)

// Decompressor is the interface that wraps the Decompress method.
type Decompressor interface {
	Decompress(to *bytes.Buffer, from io.Reader) error
}

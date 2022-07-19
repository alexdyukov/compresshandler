package decompressors

import (
	"bytes"
	"io"
)

type Decompressor interface {
	Decompress(to *bytes.Buffer, from io.Reader) error
}

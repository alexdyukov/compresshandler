package compressors

import (
	"bytes"
	"io"
)

type Compressor interface {
	Compress(level int, to io.Writer, from *bytes.Buffer) error
}

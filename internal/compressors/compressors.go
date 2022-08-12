package compressors

import (
	"bytes"
	"io"
)

type Compressor interface {
	Compress(to io.Writer, from *bytes.Buffer) error
}

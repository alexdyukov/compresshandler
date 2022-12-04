package compressor

import (
	"io"
)

type Compressor interface {
	Compress(to io.Writer, from []byte) error
}

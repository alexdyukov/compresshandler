package compressor

import (
	"io"
)

// Compressor is the interface that wraps the Compress method.
type Compressor interface {
	Compress(to io.Writer, from []byte) error
}

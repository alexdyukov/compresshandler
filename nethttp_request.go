package compresshandler

import (
	"fmt"
	"io"
	"net/http"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/gzip"
	"github.com/klauspost/compress/zlib"
)

func wrapNetHTTPRequest(req *http.Request) (*http.Request, error) {
	var (
		err             error
		contentEncoding = req.Header.Get("Content-Encoding")
	)

	for pos := 0; pos < len(contentEncoding); pos++ {
		switch contentEncoding[pos] {
		case 'z', 'Z': // z stands only in gzip >> gzip
			if req.Body, err = gzip.NewReader(req.Body); err != nil {
				return req, fmt.Errorf("compresshandler: request: failed to initialize gzip reader: %w", err)
			}
		case 'f', 'F': // f stands only in deflate >> zlib
			if req.Body, err = zlib.NewReader(req.Body); err != nil {
				return req, fmt.Errorf("compresshandler: request: failed to initialize zlib reader: %w", err)
			}
		case 'b', 'B': // b stands only in br >> brotli
			req.Body = io.NopCloser(brotli.NewReader(req.Body))
		case 'c', 'C': // c stands only in compress >> lzw
			return req, ErrNotSupported
		}
	}

	return req, nil
}

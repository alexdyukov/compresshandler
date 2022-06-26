package compresshandler

import (
	"fmt"
	"io"
	"net/http"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/gzip"
	"github.com/klauspost/compress/zlib"
)

type unsupportedContentEncodingError string

func (err unsupportedContentEncodingError) Error() string {
	return "unsupported content encoding: " + string(err)
}

func wrapNetHTTPRequest(req *http.Request) (*http.Request, error) {
	var (
		err                   error
		contentEncodingHeader = req.Header.Get("Content-Encoding")
		contentEncoding       = []byte(contentEncodingHeader)
	)

	for contentType, pos := 0, 0; pos < len(contentEncoding); pos++ {
		contentType, pos = getEncodingType(contentEncoding, pos)
		switch contentType {
		case brType:
			req.Body = io.NopCloser(brotli.NewReader(req.Body))
		case gzipType:
			if req.Body, err = gzip.NewReader(req.Body); err != nil {
				return req, fmt.Errorf("failed to initialize gzip reader: %w", err)
			}
		case deflateType:
			if req.Body, err = zlib.NewReader(req.Body); err != nil {
				return req, fmt.Errorf("failed to initialize zlib reader: %w", err)
			}
		case identityType:
			continue
		default:
			return req, unsupportedContentEncodingError(contentEncodingHeader)
		}
	}

	return req, nil
}

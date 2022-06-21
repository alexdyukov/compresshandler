package compresshandler

import (
	"io"
	"net/http"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/gzip"
	"github.com/klauspost/compress/zlib"
)

func wrapNetHTTPRequest(r *http.Request) (*http.Request, error) {
	contentEncoding := r.Header.Get("Content-Encoding")

	requestCompression := getRequestCompression([]byte(contentEncoding))
	for i := 0; i < len(requestCompression); i += 1 {
		switch requestCompression[i] {
		case gzipType:
			gzipReader, err := gzip.NewReader(r.Body)
			if err != nil {
				return r, err
			}
			r.Body = gzipReader
		case zlibType:
			zlibReader, err := zlib.NewReader(r.Body)
			if err != nil {
				return r, err
			}
			r.Body = zlibReader
		case brotliType:
			brotliReader := brotli.NewReader(r.Body)
			r.Body = io.NopCloser(brotliReader)
		case lzwType:
			return r, ErrNotSupported
		}
	}

	return r, nil
}

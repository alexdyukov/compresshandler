package compresshandler

import (
	"net/http"

	"github.com/klauspost/compress/gzip"
	"github.com/klauspost/compress/zlib"
)

type Config struct {
	GzipLevel        int
	ZlibLevel        int
	MinContentLength int
}

func NewNetHTTPHandler(config Config) func(next http.Handler) http.Handler {
	if config.GzipLevel < gzip.BestSpeed || config.GzipLevel > gzip.BestCompression {
		config.GzipLevel = gzip.DefaultCompression
	}

	if config.ZlibLevel < zlib.BestSpeed || config.ZlibLevel > zlib.BestCompression {
		config.ZlibLevel = zlib.DefaultCompression
	}

	if config.MinContentLength < 0 {
		config.MinContentLength = 0
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wrappedRequest, err := wrapNetHTTPRequest(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			acceptEncoding := wrappedRequest.Header.Get("Accept-Encoding")
			wrappedResponse := wrapNetHTTPResponse(w, config, acceptEncoding)

			next.ServeHTTP(wrappedResponse, wrappedRequest)

			if err = wrappedResponse.Flush(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})
	}
}

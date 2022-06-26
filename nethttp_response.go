package compresshandler

import (
	"fmt"
	"net/http"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/gzip"
	"github.com/klauspost/compress/zlib"
)

type wrappedResponseWriter struct {
	httpResponseWriter   http.ResponseWriter
	acceptEncoding       []byte
	bufferedResponseBody []byte
	statusCode           int
	config               Config
}

func wrapNetHTTPResponseWriter(w http.ResponseWriter, c Config, acceptEncoding string) *wrappedResponseWriter {
	return &wrappedResponseWriter{
		httpResponseWriter:   w,
		config:               c,
		acceptEncoding:       []byte(acceptEncoding),
		statusCode:           http.StatusOK,
		bufferedResponseBody: nil,
	}
}

func (wrw *wrappedResponseWriter) Header() http.Header {
	return wrw.httpResponseWriter.Header()
}

func (wrw *wrappedResponseWriter) Write(a []byte) (int, error) {
	wrw.bufferedResponseBody = append(wrw.bufferedResponseBody, a...)

	return len(a), nil
}

func (wrw *wrappedResponseWriter) WriteHeader(statusCode int) {
	wrw.statusCode = statusCode
}

func (wrw *wrappedResponseWriter) Flush() error {
	if wrw.Header().Get("Content-Type") == "" {
		wrw.Header().Set("Content-Type", http.DetectContentType(wrw.bufferedResponseBody))
	}

	if len(wrw.bufferedResponseBody) < wrw.config.MinContentLength {
		return wrw.flushNoCompression()
	}

	switch getPreferedCompression(wrw.acceptEncoding) {
	case gzipType:
		return wrw.flushGzipCompression()
	case deflateType:
		return wrw.flushZlibCompression()
	case brType:
		return wrw.flushBrotliCompression()
	default:
		return wrw.flushNoCompression()
	}
}

func (wrw *wrappedResponseWriter) flushNoCompression() error {
	wrw.httpResponseWriter.WriteHeader(wrw.statusCode)

	if _, err := wrw.httpResponseWriter.Write(wrw.bufferedResponseBody); err != nil {
		return fmt.Errorf("failed to flush uncompressed data: %w", err)
	}

	return nil
}

func (wrw *wrappedResponseWriter) flushGzipCompression() error {
	gzipWriter, err := gzip.NewWriterLevel(wrw.httpResponseWriter, wrw.config.GzipLevel)
	if err != nil { // unreacheable code, cause it raises on wrong compress level
		return fmt.Errorf("failed to initialize gzip writer: %w", err)
	}

	wrw.httpResponseWriter.Header().Set("Content-Encoding", "gzip")
	wrw.httpResponseWriter.WriteHeader(wrw.statusCode)

	if _, err = gzipWriter.Write(wrw.bufferedResponseBody); err != nil {
		return fmt.Errorf("failed to write gziped data: %w", err)
	}

	if err = gzipWriter.Close(); err != nil {
		return fmt.Errorf("failed to flush gziped data: %w", err)
	}

	return nil
}

func (wrw *wrappedResponseWriter) flushZlibCompression() error {
	zlibWriter, err := zlib.NewWriterLevel(wrw.httpResponseWriter, wrw.config.ZlibLevel)
	if err != nil { // unreacheable code, cause it raises on wrong compress level
		return fmt.Errorf("failed to initialize zlib writer: %w", err)
	}

	wrw.httpResponseWriter.Header().Set("Content-Encoding", "deflate")
	wrw.httpResponseWriter.WriteHeader(wrw.statusCode)

	if _, err = zlibWriter.Write(wrw.bufferedResponseBody); err != nil {
		return fmt.Errorf("failed to write zlibed data: %w", err)
	}

	if err = zlibWriter.Close(); err != nil {
		return fmt.Errorf("failed to flush zlibed data: %w", err)
	}

	return nil
}

func (wrw *wrappedResponseWriter) flushBrotliCompression() error {
	brotliWriter := brotli.NewWriterLevel(wrw.httpResponseWriter, wrw.config.BrotliLevel)

	wrw.httpResponseWriter.Header().Set("Content-Encoding", "br")
	wrw.httpResponseWriter.WriteHeader(wrw.statusCode)

	if _, err := brotliWriter.Write(wrw.bufferedResponseBody); err != nil {
		return fmt.Errorf("failed to write brotlied data: %w", err)
	}

	if err := brotliWriter.Close(); err != nil {
		return fmt.Errorf("failed to flush brotlied data: %w", err)
	}

	return nil
}

// Package compresshandler provides http middlewares which auto compress http responses and
// auto decompress http requests based on request's Accept-Encoding and Content-Encoding headers.
package compresshandler

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/alexdyukov/compresshandler/internal/compressor"
	"github.com/alexdyukov/compresshandler/internal/decompressor"
	"github.com/alexdyukov/compresshandler/internal/encoding"
)

type (
	compressors   map[int]compressor.Compressor
	decompressors map[int]decompressor.Decompressor
	// Config describes configuration for New* methods.
	Config struct {
		GzipLevel        int // gzip compression level for response between GzipBestSpeed and GzipBestCompression including boundaries
		ZlibLevel        int // zlib compression level for response between ZlibBestSpeed and ZlibBestCompression including boundaries
		BrotliLevel      int // brotli compression level for response between BrotliBestSpeed and BrotliBestCompression including boundaries
		MinContentLength int // minimal length of raw response to be compress
	}
	wrappedNetHTTPResponseWriter struct {
		http.ResponseWriter
		bufferedResponse *bytes.Buffer
		statusCode       *int
	}
)

const (
	GzipBestSpeed            = compressor.GzipBestSpeed
	GzipBestCompression      = compressor.GzipBestCompression
	GzipDefaultCompression   = (GzipBestCompression - GzipBestSpeed) / 2
	ZlibBestSpeed            = compressor.ZlibBestSpeed
	ZlibBestCompression      = compressor.ZlibBestCompression
	ZlibDefaultCompression   = (ZlibBestCompression - ZlibBestSpeed) / 2
	BrotliBestSpeed          = compressor.BrotliBestSpeed
	BrotliBestCompression    = compressor.BrotliBestCompression
	BrotliDefaultCompression = (BrotliBestCompression - BrotliBestSpeed) / 2
)

func (cfg *Config) fix() {
	if cfg.GzipLevel < GzipBestSpeed || cfg.GzipLevel > GzipBestCompression {
		cfg.GzipLevel = GzipDefaultCompression
	}

	if cfg.ZlibLevel < ZlibBestSpeed || cfg.ZlibLevel > ZlibBestCompression {
		cfg.ZlibLevel = ZlibDefaultCompression
	}

	if cfg.BrotliLevel < BrotliBestSpeed || cfg.BrotliLevel > BrotliBestCompression {
		cfg.BrotliLevel = BrotliDefaultCompression
	}

	if cfg.MinContentLength < 0 {
		cfg.MinContentLength = 0
	}
}

func (cfg *Config) getPossibleCompressors() compressors {
	cfg.fix()

	return compressors{
		encoding.BrType:      compressor.NewBrotli(cfg.BrotliLevel),
		encoding.GzipType:    compressor.NewGzip(cfg.GzipLevel),
		encoding.DeflateType: compressor.NewZlib(cfg.ZlibLevel),
	}
}

func (cfg *Config) getPossibleDecompressors() decompressors {
	return decompressors{
		encoding.BrType:      decompressor.NewBrotli(),
		encoding.GzipType:    decompressor.NewGzip(),
		encoding.DeflateType: decompressor.NewZlib(),
	}
}

func (wrapped *wrappedNetHTTPResponseWriter) Write(a []byte) (int, error) {
	n, err := wrapped.bufferedResponse.Write(a)

	return n, fmt.Errorf("%w", err)
}

func (wrapped *wrappedNetHTTPResponseWriter) WriteHeader(statusCode int) {
	*(wrapped.statusCode) = statusCode
}

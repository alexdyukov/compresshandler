package compresshandler

import (
	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/gzip"
	"github.com/klauspost/compress/zlib"
)

type Config struct {
	GzipLevel        int
	ZlibLevel        int
	BrotliLevel      int
	MinContentLength int
}

func (cfg *Config) fix() {
	if cfg.GzipLevel < gzip.BestSpeed || cfg.GzipLevel > gzip.BestCompression {
		cfg.GzipLevel = gzip.DefaultCompression
	}

	if cfg.ZlibLevel < zlib.BestSpeed || cfg.ZlibLevel > zlib.BestCompression {
		cfg.ZlibLevel = zlib.DefaultCompression
	}

	if cfg.BrotliLevel < brotli.BestSpeed || cfg.BrotliLevel > brotli.BestCompression {
		cfg.BrotliLevel = brotli.DefaultCompression
	}

	if cfg.MinContentLength < 0 {
		cfg.MinContentLength = 0
	}
}

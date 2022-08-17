package compresshandler_test

import (
	"net/http"

	"github.com/alexdyukov/compresshandler/internal/compressor"
	"github.com/alexdyukov/compresshandler/internal/decompressor"
	"github.com/alexdyukov/compresshandler/internal/encoding"
	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/gzip"
	"github.com/klauspost/compress/zlib"
	"github.com/valyala/fasthttp"
)

var (
	comps map[int]compressor.Compressor = map[int]compressor.Compressor{
		encoding.BrType:      compressor.NewBrotli((brotli.BestCompression - brotli.BestSpeed) / 2),
		encoding.GzipType:    compressor.NewGzip((gzip.BestCompression - gzip.BestSpeed) / 2),
		encoding.DeflateType: compressor.NewZlib((zlib.BestCompression - zlib.BestSpeed) / 2),
	}
	decomps map[int]decompressor.Decompressor = map[int]decompressor.Decompressor{
		encoding.BrType:      decompressor.NewBrotli(),
		encoding.GzipType:    decompressor.NewGzip(),
		encoding.DeflateType: decompressor.NewZlib(),
	}
	testString               = "there is A test string !@#$%^&*()_+"
	netHTTPReturnStatusCode  = http.StatusAccepted
	fastHTTPReturnStatusCode = fasthttp.StatusAccepted
	tests                    = []struct {
		name           string
		acceptEncoding string
		requestType    int
		responseType   int
	}{
		{
			name:           "vanilla request vanilla response",
			acceptEncoding: " ",
			requestType:    encoding.IdentityType,
			responseType:   encoding.IdentityType,
		}, {
			name:           "vanilla request gzip response",
			acceptEncoding: "gzip",
			requestType:    encoding.IdentityType,
			responseType:   encoding.GzipType,
		}, {
			name:           "gzip request vanilla response",
			acceptEncoding: "identity",
			requestType:    encoding.GzipType,
			responseType:   encoding.IdentityType,
		}, {
			name:           "gzip request zlib response",
			acceptEncoding: "deflate",
			requestType:    encoding.GzipType,
			responseType:   encoding.DeflateType,
		},
	}
)

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

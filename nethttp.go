package compresshandler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/alexdyukov/compresshandler/internal/compressors"
	"github.com/alexdyukov/compresshandler/internal/decompressors"
	"github.com/alexdyukov/compresshandler/internal/encoding"
)

type (
	availableCompressors         map[int]compressors.Compressor
	availableDecompressors       map[int]decompressors.Decompressor
	wrappedNetHTTPResponseWriter struct {
		http.ResponseWriter
		bufferedResponse *bytes.Buffer
		statusCode       *int
	}
)

func (wrw *wrappedNetHTTPResponseWriter) Write(a []byte) (int, error) {
	n, err := wrw.bufferedResponse.Write(a)

	return n, fmt.Errorf("%w", err)
}

func (wrw *wrappedNetHTTPResponseWriter) WriteHeader(statusCode int) {
	*(wrw.statusCode) = statusCode
}

func NewNetHTTP(config Config) func(next http.Handler) http.Handler {
	config.fix()

	bufferPool := &sync.Pool{
		New: func() interface{} {
			return &bytes.Buffer{}
		},
	}

	comps := availableCompressors{
		encoding.BrType:      compressors.NewBrotli(config.BrotliLevel),
		encoding.GzipType:    compressors.NewGzip(config.GzipLevel),
		encoding.DeflateType: compressors.NewZlib(config.ZlibLevel),
	}

	decomps := availableDecompressors{
		encoding.BrType:      decompressors.NewBrotli(),
		encoding.GzipType:    decompressors.NewGzip(),
		encoding.DeflateType: decompressors.NewZlib(),
	}

	return func(next http.Handler) http.Handler {
		return decompressNetHTTP(bufferPool,
			decomps,
			compressNetHTTP(config.MinContentLength, bufferPool, comps, next),
		)
	}
}

func decompressNetHTTP(buffers *sync.Pool, decomps availableDecompressors, next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		if request.Header.Get("Upgrade") != "" {
			next.ServeHTTP(responseWriter, request)

			return
		}

		encodings, err := encoding.ParseContentEncoding([]byte(request.Header.Get("Content-Encoding")))
		if err != nil {
			http.Error(responseWriter, "not supported Content-Encoding", http.StatusBadRequest)

			return
		}

		usedBuffers := []*bytes.Buffer{}
		defer func() {
			for _, v := range usedBuffers {
				buffers.Put(v)
			}
		}()

		for enc := 0; enc < len(encodings); enc++ {
			buffer, okay := buffers.Get().(*bytes.Buffer)
			if !okay {
				panic("unreachable code")
			}
			buffer.Reset()
			usedBuffers = append(usedBuffers, buffer)

			decompressor := decomps[encodings[enc]]
			if err = decompressor.Decompress(buffer, request.Body); err != nil {
				http.Error(responseWriter, "invalid request body with presented Content-Encoding", http.StatusBadRequest)

				return
			}

			request.Body = io.NopCloser(buffer)
		}

		next.ServeHTTP(responseWriter, request)
	})
}

func compressNetHTTP(minLength int, buffers *sync.Pool, comps availableCompressors, next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		statusCode := http.StatusOK

		buffer, okay := buffers.Get().(*bytes.Buffer)
		if !okay {
			panic("unreachable code")
		}
		buffer.Reset()
		defer buffers.Put(buffer)

		next.ServeHTTP(&wrappedNetHTTPResponseWriter{
			ResponseWriter:   responseWriter,
			bufferedResponse: buffer,
			statusCode:       &statusCode,
		}, request)

		if responseWriter.Header().Get("Content-Type") == "" {
			responseWriter.Header().Set("Content-Type", http.DetectContentType(buffer.Bytes()))
		}

		preferedEncoding := encoding.ParseAcceptEncoding([]byte(request.Header.Get("Accept-Encoding")))

		compressor, okay := comps[preferedEncoding]
		if !okay || buffer.Len() < minLength {
			responseWriter.WriteHeader(statusCode)
			_, err := responseWriter.Write(buffer.Bytes())
			if err != nil {
				panic(err)
			}

			return
		}

		responseWriter.Header().Set("Content-Encoding", encoding.ToString(preferedEncoding))
		responseWriter.WriteHeader(statusCode)

		err := compressor.Compress(responseWriter, buffer)
		if err != nil {
			panic(err)
		}
	})
}

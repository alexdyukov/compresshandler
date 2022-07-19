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

type wrappedNetHTTPResponseWriter struct {
	http.ResponseWriter
	bufferedResponse *bytes.Buffer
	statusCode       *int
}

func (wrw *wrappedNetHTTPResponseWriter) Write(a []byte) (int, error) {
	n, err := wrw.bufferedResponse.Write(a)

	return n, fmt.Errorf("%w", err)
}

func (wrw *wrappedNetHTTPResponseWriter) WriteHeader(statusCode int) {
	*(wrw.statusCode) = statusCode
}

func NewNetHTTP(config Config) func(next http.Handler) http.Handler {
	config.fix()

	var (
		levels  = config.getLevels()
		buffers = &sync.Pool{
			New: func() interface{} {
				return &bytes.Buffer{}
			},
		}
		availableCompressors = map[int]compressors.Compressor{
			encoding.BrType:      compressors.NewBrotli(),
			encoding.GzipType:    compressors.NewGzip(),
			encoding.DeflateType: compressors.NewZlib(),
		}
		availableDecompressors = map[int]decompressors.Decompressor{
			encoding.BrType:      decompressors.NewBrotli(),
			encoding.GzipType:    decompressors.NewGzip(),
			encoding.DeflateType: decompressors.NewZlib(),
		}
	)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
			var (
				usedBuffers            []*bytes.Buffer
				parsedContentEncodings []int
				buffer                 *bytes.Buffer
				compressor             compressors.Compressor
				okay                   bool
			)
			defer func() {
				for _, v := range usedBuffers {
					buffers.Put(v)
				}
			}()

			if request.Header.Get("Upgrade") != "" {
				next.ServeHTTP(responseWriter, request)

				return
			}

			parsedContentEncodings, err := encoding.ParseContentEncoding([]byte(request.Header.Get("Content-Encoding")))
			if err != nil {
				http.Error(responseWriter, "not implemented encoding in Content-Encoding header", http.StatusBadRequest)

				return
			}

			for enc := 0; enc < len(parsedContentEncodings); enc++ {
				buffer, okay = buffers.Get().(*bytes.Buffer)
				if !okay {
					panic("unreachable code")
				}
				buffer.Reset()
				usedBuffers = append(usedBuffers, buffer)

				decompressor := availableDecompressors[parsedContentEncodings[enc]]
				if err = decompressor.Decompress(buffer, request.Body); err != nil {
					http.Error(responseWriter, "invalid request body with presented Content-Encoding", http.StatusBadRequest)

					return
				}

				request.Body = io.NopCloser(buffer)
			}

			statusCode := http.StatusOK

			buffer, okay = buffers.Get().(*bytes.Buffer)
			if !okay {
				panic("unreachable code")
			}
			buffer.Reset()

			usedBuffers = append(usedBuffers, buffer)

			next.ServeHTTP(&wrappedNetHTTPResponseWriter{
				ResponseWriter:   responseWriter,
				bufferedResponse: buffer,
				statusCode:       &statusCode,
			}, request)

			if responseWriter.Header().Get("Content-Type") == "" {
				responseWriter.Header().Set("Content-Type", http.DetectContentType(buffer.Bytes()))
			}

			preferedEncoding := encoding.ParseAcceptEncoding([]byte(request.Header.Get("Accept-Encoding")))

			compressor, okay = availableCompressors[preferedEncoding]
			if !okay || buffer.Len() < config.MinContentLength {
				responseWriter.WriteHeader(statusCode)
				_, err = responseWriter.Write(buffer.Bytes())
				if err != nil {
					panic(err)
				}

				return
			}

			responseWriter.Header().Set("Content-Encoding", encoding.ToString(preferedEncoding))
			responseWriter.WriteHeader(statusCode)

			err = compressor.Compress(levels[preferedEncoding], responseWriter, buffer)
			if err != nil {
				panic(err)
			}
		})
	}
}

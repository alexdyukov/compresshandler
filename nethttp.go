package compresshandler

import (
	"bytes"
	"io"
	"net/http"
	"sync"

	"github.com/alexdyukov/compresshandler/internal/encoding"
)

func NewNetHTTP(config Config) func(next http.Handler) http.Handler {
	bufferPool := &sync.Pool{
		New: func() interface{} {
			return &bytes.Buffer{}
		},
	}

	comps := config.getPossibleCompressors()

	decomps := config.getPossibleDecompressors()

	minlen := config.MinContentLength

	return func(next http.Handler) http.Handler {
		wrapped := decompressNetHTTP(bufferPool, decomps, next)
		wrapped = compressNetHTTP(minlen, bufferPool, comps, wrapped)

		return wrapped
	}
}

func decompressNetHTTP(bufferPool *sync.Pool, decomps decompressors, next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		encodings, err := encoding.ParseContentEncoding([]byte(request.Header.Get("Content-Encoding")))
		if err != nil {
			http.Error(responseWriter, "not supported Content-Encoding", http.StatusBadRequest)

			return
		}

		usedBuffers := []*bytes.Buffer{}
		defer func() {
			for _, v := range usedBuffers {
				bufferPool.Put(v)
			}
		}()

		for enc := 0; enc < len(encodings); enc++ {
			buffer, okay := bufferPool.Get().(*bytes.Buffer)
			if !okay {
				panic("unreachable code")
			}
			buffer.Reset()
			usedBuffers = append(usedBuffers, buffer)

			decomp := decomps[encodings[enc]]
			if err = decomp.Decompress(buffer, request.Body); err != nil {
				http.Error(responseWriter, "invalid request body with presented Content-Encoding", http.StatusBadRequest)

				return
			}

			request.Body = io.NopCloser(buffer)
		}

		request.Header.Del("Content-Encoding")

		next.ServeHTTP(responseWriter, request)
	})
}

func compressNetHTTP(minLength int, bufferPool *sync.Pool, comps compressors, next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		statusCode := http.StatusOK

		upstreamResponse, okay := bufferPool.Get().(*bytes.Buffer)
		if !okay {
			panic("unreachable code")
		}
		upstreamResponse.Reset()
		defer bufferPool.Put(upstreamResponse)

		next.ServeHTTP(&wrappedNetHTTPResponseWriter{
			ResponseWriter:   responseWriter,
			bufferedResponse: upstreamResponse,
			statusCode:       &statusCode,
		}, request)

		upstreamResponseBody := upstreamResponse.Bytes()

		if responseWriter.Header().Get("Content-Type") == "" {
			responseWriter.Header().Set("Content-Type", http.DetectContentType(upstreamResponseBody))
		}

		preferedEncoding := encoding.ParseAcceptEncoding([]byte(request.Header.Get("Accept-Encoding")))

		comp, okay := comps[preferedEncoding]
		if !okay || upstreamResponse.Len() < minLength || responseWriter.Header().Get("Content-Encoding") != "" {
			responseWriter.WriteHeader(statusCode)
			_, err := responseWriter.Write(upstreamResponse.Bytes())
			if err != nil {
				panic(err)
			}

			return
		}

		responseWriter.Header().Set("Content-Encoding", encoding.ToString(preferedEncoding))
		responseWriter.WriteHeader(statusCode)

		err := comp.Compress(responseWriter, upstreamResponseBody)
		if err != nil {
			panic(err)
		}
	})
}

package compresshandler_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alexdyukov/compresshandler"
	"github.com/alexdyukov/compresshandler/internal/compressors"
	"github.com/alexdyukov/compresshandler/internal/decompressors"
	"github.com/alexdyukov/compresshandler/internal/encoding"
	"github.com/stretchr/testify/assert"
)

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func TestNetHttp(t *testing.T) {
	var (
		availableCompressors map[int]compressors.Compressor = map[int]compressors.Compressor{
			encoding.BrType:      compressors.NewBrotli(),
			encoding.GzipType:    compressors.NewGzip(),
			encoding.DeflateType: compressors.NewZlib(),
		}
		availableDecompressors map[int]decompressors.Decompressor = map[int]decompressors.Decompressor{
			encoding.BrType:      decompressors.NewBrotli(),
			encoding.GzipType:    decompressors.NewGzip(),
			encoding.DeflateType: decompressors.NewZlib(),
		}
		testString           = "there is A test string !@#$%^&*()_+"
		httpReturnStatusCode = http.StatusAccepted
		testCompressLevel    = 6
	)

	tests := []struct {
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

	compress := compresshandler.NewNetHTTP(compresshandler.Config{})
	revHTTP := compress(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		input := strings.TrimSpace(string(b))
		output := reverse(input)
		w.WriteHeader(httpReturnStatusCode)
		w.Write([]byte(output))
	}))

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			requestBody := &bytes.Buffer{}

			reader := bytes.NewBufferString(testString)
			compressor, needCompress := availableCompressors[test.requestType]
			if !needCompress {
				requestBody.ReadFrom(reader)
			} else {
				err := compressor.Compress(testCompressLevel, requestBody, reader)
				if err != nil {
					t.Fatalf("cannot Compress requestBody with encodingType (%v): %v", test.requestType, err)
				}
			}

			request := httptest.NewRequest(http.MethodPost, "/", requestBody)

			switch test.requestType {
			case encoding.GzipType:
				request.Header.Set("Content-Encoding", "gzip")
			case encoding.DeflateType:
				request.Header.Set("Content-Encoding", "deflate")
			case encoding.BrType:
				request.Header.Set("Content-Encoding", "br")
			}

			request.Header.Set("Accept-Encoding", test.acceptEncoding)

			recorder := httptest.NewRecorder()
			revHTTP.ServeHTTP(recorder, request)
			response := recorder.Result()
			defer response.Body.Close()

			assert.Equal(t, httpReturnStatusCode, response.StatusCode)

			cleanedResponseBody := &bytes.Buffer{}

			decompressor, needDecompress := availableDecompressors[test.responseType]
			if !needDecompress {
				cleanedResponseBody.ReadFrom(response.Body)
			} else {
				err := decompressor.Decompress(cleanedResponseBody, response.Body)
				if err != nil {
					t.Fatalf("cannot Decompress response.Body with encodingType (%v): %v", test.responseType, err)
				}
			}

			responseContentEncoding := response.Header.Get("Content-Encoding")
			switch test.responseType {
			case encoding.GzipType:
				assert.Contains(t, responseContentEncoding, "gzip")
			case encoding.DeflateType:
				assert.Contains(t, responseContentEncoding, "deflate")
			case encoding.BrType:
				assert.Contains(t, responseContentEncoding, "br")
			}

			assert.True(t, strings.TrimSpace(string(cleanedResponseBody.Bytes())) == reverse(testString))
		})
	}
}

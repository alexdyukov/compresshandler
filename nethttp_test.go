package compresshandler_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alexdyukov/compresshandler"
	"github.com/alexdyukov/compresshandler/internal/encoding"
	"github.com/stretchr/testify/assert"
)

func TestNetHttp(t *testing.T) {
	compress := compresshandler.NewNetHTTP(compresshandler.Config{})
	revNetHTTP := compress(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		input := strings.TrimSpace(string(b))
		output := reverse(input)
		w.WriteHeader(netHTTPReturnStatusCode)
		w.Write([]byte(output))
	}))

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			requestBody := &bytes.Buffer{}

			comp, needCompress := comps[test.requestType]
			if !needCompress {
				requestBody.ReadFrom(bytes.NewBufferString(testString))
			} else {
				err := comp.Compress(requestBody, []byte(testString))
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
			revNetHTTP.ServeHTTP(recorder, request)
			response := recorder.Result()
			defer response.Body.Close()

			assert.Equal(t, netHTTPReturnStatusCode, response.StatusCode)

			cleanedResponseBody := &bytes.Buffer{}

			decomp, needDecompress := decomps[test.responseType]
			if !needDecompress {
				cleanedResponseBody.ReadFrom(response.Body)
			} else {
				err := decomp.Decompress(cleanedResponseBody, response.Body)
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

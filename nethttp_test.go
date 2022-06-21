package compresshandler

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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
	testString := "there is A test string !@#$%^&*()_+"

	tests := []struct {
		name           string
		acceptEncoding string
		requestType    int
		responseType   int
	}{
		{
			name:           "vanilla request vanilla response",
			acceptEncoding: " ",
			requestType:    identityType,
			responseType:   identityType,
		}, {
			name:           "vanilla request gzip response",
			acceptEncoding: "gzip",
			requestType:    identityType,
			responseType:   gzipType,
		}, {
			name:           "gzip request vanilla response",
			acceptEncoding: "identity",
			requestType:    gzipType,
			responseType:   identityType,
		}, {
			name:           "gzip request zlib response",
			acceptEncoding: "deflate",
			requestType:    gzipType,
			responseType:   zlibType,
		},
	}

	compressHandler := NewNetHTTPHandler(Config{})
	revHTTP := compressHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		input := strings.TrimSpace(string(b))
		output := reverse(input)
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(output))
	}))

	// run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//request
			var request *http.Request

			switch tt.requestType {
			case gzipType:
				requestBody, err := gzipSlice([]byte(testString))
				if err != nil {
					panic("cannot initialize requestBody with gzipType")
				}
				request = httptest.NewRequest(http.MethodGet, "/", bytes.NewReader(requestBody))
				request.Header.Set("Content-Encoding", "gzip")
			case zlibType:
				requestBody, err := zlibSlice([]byte(testString))
				if err != nil {
					panic("cannot initialize requestBody with zlibType")
				}
				request = httptest.NewRequest(http.MethodGet, "/", bytes.NewReader(requestBody))
				request.Header.Set("Content-Encoding", "deflate")
			default: // identityType
				requestBody := []byte(testString)
				request = httptest.NewRequest(http.MethodGet, "/", bytes.NewReader(requestBody))
			}

			switch tt.responseType {
			case gzipType:
				request.Header.Set("Accept-Encoding", "gzip")
			case zlibType:
				request.Header.Set("Accept-Encoding", "deflate")
			default:
				request.Header.Set("Accept-Encoding", "identity")
			}

			//response
			recorder := httptest.NewRecorder()
			revHTTP.ServeHTTP(recorder, request)
			response := recorder.Result()
			defer response.Body.Close()

			//checks
			assert.Equal(t, http.StatusAccepted, response.StatusCode)

			returnedBody, err := io.ReadAll(response.Body)
			assert.Nil(t, err)

			var uncompressedReturnBody []byte

			switch tt.responseType {
			case gzipType:
				assert.Contains(t, response.Header.Get("Content-Encoding"), "gzip")
				uncompressedReturnBody, err = ungzipSlice(returnedBody)
				if err != nil {
					panic("cannot get uncompressed body from gzipped response")
				}
			case zlibType:
				assert.Contains(t, response.Header.Get("Content-Encoding"), "deflate")
				uncompressedReturnBody, err = unzlibSlice(returnedBody)
				if err != nil {
					panic("cannot get uncompressed body from zlibbed response")
				}
			default:
				uncompressedReturnBody = returnedBody
			}

			assert.True(t, strings.TrimSpace(string(uncompressedReturnBody)) == reverse(testString))
		})
	}
}

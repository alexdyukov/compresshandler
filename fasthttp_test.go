package compresshandler_test

import (
	"bytes"
	"net"
	"strings"
	"testing"

	"github.com/alexdyukov/compresshandler"
	"github.com/alexdyukov/compresshandler/internal/encoding"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
)

func TestFastHttp(t *testing.T) {
	compress := compresshandler.NewFastHTTP(compresshandler.Config{})
	revFastHTTP := compress(func(ctx *fasthttp.RequestCtx) {
		input := strings.TrimSpace(string(ctx.Request.Body()))
		output := reverse(input)
		ctx.SetStatusCode(fastHTTPReturnStatusCode)
		ctx.SetBodyString(output)
	})

	listener := fasthttputil.NewInmemoryListener()
	defer listener.Close()

	go func() {
		err := fasthttp.Serve(listener, revFastHTTP)
		if err != nil {
			panic(err)
		}
	}()

	httpClient := fasthttp.Client{
		Dial: func(addr string) (net.Conn, error) {
			return listener.Dial()
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
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

			request := fasthttp.AcquireRequest()
			request.SetHost("test")
			request.Header.SetMethod(fasthttp.MethodPost)
			request.SetBody(requestBody.Bytes())

			switch test.requestType {
			case encoding.GzipType:
				request.Header.Set("Content-Encoding", "gzip")
			case encoding.DeflateType:
				request.Header.Set("Content-Encoding", "deflate")
			case encoding.BrType:
				request.Header.Set("Content-Encoding", "br")
			}

			request.Header.Set("Accept-Encoding", test.acceptEncoding)

			response := fasthttp.AcquireResponse()
			defer fasthttp.ReleaseResponse(response)

			err := httpClient.Do(request, response)
			if err != nil {
				t.Fatalf("cannot make http request to compress handler: %v", err)
			}

			assert.Equal(t, fastHTTPReturnStatusCode, response.StatusCode())

			cleanedResponseBody := &bytes.Buffer{}
			responseBody := bytes.NewReader(response.Body())

			decomp, needDecompress := decomps[test.responseType]
			if !needDecompress {
				cleanedResponseBody.ReadFrom(responseBody)
			} else {
				err := decomp.Decompress(cleanedResponseBody, responseBody)
				if err != nil {
					t.Fatalf("cannot Decompress response.Body with encodingType (%v): %v", test.responseType, err)
				}
			}

			responseContentEncoding := string(response.Header.Peek("Content-Encoding"))
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

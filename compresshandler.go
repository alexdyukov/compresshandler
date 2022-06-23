package compresshandler

import (
	"errors"
	"net/http"
)

var ErrNotSupported = errors.New("unsupported content encoding")

func NewNetHTTPHandler(config Config) func(next http.Handler) http.Handler {
	config.fix()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
			wrappedRequest, err := wrapNetHTTPRequest(request)
			if err != nil {
				http.Error(responseWriter, err.Error(), http.StatusBadRequest)

				return
			}

			acceptEncoding := wrappedRequest.Header.Get("Accept-Encoding")
			wrappedResponseWriter := wrapNetHTTPResponseWriter(responseWriter, config, acceptEncoding)

			next.ServeHTTP(wrappedResponseWriter, wrappedRequest)

			if err = wrappedResponseWriter.Flush(); err != nil {
				http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
			}
		})
	}
}

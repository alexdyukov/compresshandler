package compresshandler

import (
	"net/http"
)

func NewNetHTTP(config Config) func(next http.Handler) http.Handler {
	config.fix()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
			if request.Header.Get("Upgrade") != "" {
				next.ServeHTTP(responseWriter, request)

				return
			}

			wrappedRequest, err := wrapNetHTTPRequest(request)
			if err != nil {
				http.Error(responseWriter, err.Error(), http.StatusBadRequest)

				return
			}

			acceptEncoding := wrappedRequest.Header.Get("Accept-Encoding")
			wrappedRequest.Header.Del("Accept-Encoding")

			wrappedResponseWriter := wrapNetHTTPResponseWriter(responseWriter, config, acceptEncoding)

			next.ServeHTTP(wrappedResponseWriter, wrappedRequest)

			if err = wrappedResponseWriter.Flush(); err != nil {
				http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
			}
		})
	}
}

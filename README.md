# compresshandler
Go compress http handler. Just fire and forget.
====
[![GoDoc](https://godoc.org/github.com/alexdyukov/compresshandler?status.svg)](https://godoc.org/github.com/alexdyukov/compresshandler)
[![CI](https://github.com/alexdyukov/compresshandler/actions/workflows/golangci-lint.yml/badge.svg?branch=master)](https://github.com/alexdyukov/compresshandler/actions/workflows/golangci-lint.yml?query=branch%3Amaster)

Package provides methods to wrap http handler for auto decompress compressed data & compress response with prefered client compression. That is. Nothing more or nothing less.

## Example

```go
import (
        "net/http"

        "github.com/alexdyukov/compresshandler"
        "github.com/klauspost/compress/gzip"
        "github.com/klauspost/compress/zlib"
        "github.com/andybalholm/brotli"
)

func main() {
        echo := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                w.Write(r.Body)
        })

        compressConfig := compresshandler.Config{
                GzipLevel:        gzip.DefaultCompression,
                ZlibLevel:        zlib.DefaultCompression,
                BrotliLevel:      brotli.DefaultCompression,
                MinContentLength: 1400,
        }

        http.ListenAndServe(":8080", compresshandler.NewNetHTTP(echo))
}
```

## TODOs

* fasthttp and any other handler
* Configurable content response type. We dont need to zip already zipped image or archives

## License

BSD licensed. See the included LICENSE file for details.
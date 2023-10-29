# compresshandler
Go auto compress and decompress handlers for net/http and fasthttp
====
[![GoDoc](https://godoc.org/github.com/alexdyukov/compresshandler?status.svg)](https://godoc.org/github.com/alexdyukov/compresshandler)
[![CI](https://github.com/alexdyukov/compresshandler/actions/workflows/lint.yml/badge.svg?branch=master)](https://github.com/alexdyukov/compresshandler/actions/workflows/lint.yml?query=branch%3Amaster)

This package provides a middleware for net/http and fasthttp that auto decompress request body and auto compress response body with prefered client compression. Supports all [IANA's initially registred tokens](https://www.rfc-editor.org/rfc/rfc2616#section-3.5)

## Restrictions

### Server decompressor

According to RFCs there is no 'Accept-Encoding' header at server side response. It means you cannot tell clients (browsers, include headless browsers like curl/python's request) that your server accept compressed requests. But some of the backends (for example [mod_deflate](https://httpd.apache.org/docs/2.2/mod/mod_deflate.html#input)) support compressed http requests, thats why the same feature exists in this package.

### Encoding support

There is other compression algorithm: LZW and Zstd. But overall score for encoding+transfer+decoding is the same. If you really want to increase content transfer performance, [its better](https://advancedweb.hu/revisiting-webapp-performance-on-http-2/) to use [minification](github.com/tdewolff/minify) + compression (this package) + http2, rather then jumps between algos, because:
* users does not care which one we use, because only TTI (time to interactive) counts. There is no difference between 0.28sec TTI and 0.30sec TTI
* operation team does not care which one we use, because only total cost of io/cpu/ram counts. There is no win-win algo, who dramatically decrease it from 10k$ into 1k$
* other developers too lazy to enable any non gzip compression/decompression support, because time is money

## Usage

### `net/http`

```go
package main

import (
        "io"
        "net/http"

        "github.com/alexdyukov/compresshandler"
)

func main() {
        echo := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                b, _ := io.ReadAll(r.Body)
                w.Write(b)
        })

        compressConfig := compresshandler.Config{
                GzipLevel:        compresshandler.GzipDefaultCompression,
                ZlibLevel:        compresshandler.ZlibDefaultCompression,
                BrotliLevel:      compresshandler.BrotliDefaultCompression,
                MinContentLength: 1400,
        }

        compress := compresshandler.NewNetHTTP(compressConfig)

        http.ListenAndServe(":8080", compress(echo))
}
```

### `fasthttp`

```go
package main

import (
        "github.com/alexdyukov/compresshandler"
        "github.com/valyala/fasthttp"
)

func main() {
        echo := func(ctx *fasthttp.RequestCtx) {
                ctx.SetBody(ctx.Request.Body())
        }

        compressConfig := compresshandler.Config{
                GzipLevel:        compresshandler.GzipDefaultCompression,
                ZlibLevel:        compresshandler.ZlibDefaultCompression,
                BrotliLevel:      compresshandler.BrotliDefaultCompression,
                MinContentLength: 1400,
        }

        compress := compresshandler.NewFastHTTP(compressConfig)

        fasthttp.ListenAndServe(":8080", compress(echo))
}
```

## TODOs

* configurable content response type. We dont need to compress already compressed image or archives:
```
https://www.iana.org/assignments/media-types/media-types.xhtml
# default should be:
application/*json
application/*xml
application/*json
*javascript
*ecmascript
```

## License

MIT licensed. See the included LICENSE file for details.

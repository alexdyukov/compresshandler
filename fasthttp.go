package compresshandler

import (
	"bytes"
	"net/http"
	"sync"

	"github.com/alexdyukov/compresshandler/internal/encoding"
	"github.com/valyala/fasthttp"
)

func NewFastHTTP(config Config) func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	bufferPool := &sync.Pool{
		New: func() interface{} {
			return &bytes.Buffer{}
		},
	}

	comps := config.getPossibleCompressors()

	decomps := config.getPossibleDecompressors()

	minlen := config.MinContentLength

	return func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		wrapped := decompressFastHTTP(bufferPool, decomps, next)
		wrapped = compressFastHTTP(minlen, bufferPool, comps, wrapped)

		return wrapped
	}
}

func decompressFastHTTP(bufferPool *sync.Pool, decomps decompressors, next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		encodings, err := encoding.ParseContentEncoding(ctx.Request.Header.Peek("Content-Encoding"))
		if err != nil {
			ctx.Error("not supported Content-Encoding", fasthttp.StatusBadRequest)

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
			if err = decomp.Decompress(buffer, bytes.NewReader(ctx.Request.Body())); err != nil {
				ctx.Error("invalid request body with presented Content-Encoding", fasthttp.StatusBadRequest)

				return
			}

			ctx.Request.SetBodyStream(buffer, -1)
		}

		ctx.Request.Header.Del("Content-Encoding")

		next(ctx)
	})
}

func compressFastHTTP(minLength int, bufferPool *sync.Pool, comps compressors, next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		next(ctx)

		if contentType := ctx.Response.Header.Peek("Content-Type"); len(contentType) == 0 {
			ctx.Response.Header.Set("Content-Type", http.DetectContentType(ctx.Response.Body()))
		}

		preferedEncoding := encoding.ParseAcceptEncoding(ctx.Request.Header.Peek("Accept-Encoding"))

		comp, okay := comps[preferedEncoding]
		if !okay || len(ctx.Response.Body()) < minLength || len(ctx.Response.Header.Peek("Content-Encoding")) != 0 {
			return
		}

		ctx.Response.Header.Set("Content-Encoding", encoding.ToString(preferedEncoding))

		buffer, okay := bufferPool.Get().(*bytes.Buffer)
		if !okay {
			panic("unreachable code")
		}
		buffer.Reset()
		defer bufferPool.Put(buffer)

		err := comp.Compress(buffer, ctx.Response.Body())
		if err != nil {
			panic(err)
		}

		ctx.Response.SetBody(buffer.Bytes())
	})
}

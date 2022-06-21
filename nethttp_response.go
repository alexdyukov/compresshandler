package compresshandler

import (
	"net/http"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/gzip"
	"github.com/klauspost/compress/zlib"
)

type netHTTPCompressor struct {
	responseWriter http.ResponseWriter
	config         Config
	acceptEncoding []byte
	statusCode     int
	buffer         []byte
}

func (c *netHTTPCompressor) Flush() error {
	if len(c.buffer) < c.config.MinContentLength {
		//do not compress
		c.responseWriter.WriteHeader(c.statusCode)
		_, err := c.responseWriter.Write(c.buffer)
		return err
	}

	switch getPreferedCompression(c.acceptEncoding) {
	case gzipType:
		gzipWriter, err := gzip.NewWriterLevel(c.responseWriter, c.config.GzipLevel)
		if err != nil {
			return err
		}
		defer gzipWriter.Close()

		c.responseWriter.Header().Set("Content-Encoding", "gzip")
		c.responseWriter.WriteHeader(c.statusCode)
		_, err = gzipWriter.Write(c.buffer)
		return err
	case zlibType:
		zlibWriter, err := zlib.NewWriterLevel(c.responseWriter, c.config.ZlibLevel)
		if err != nil {
			return err
		}
		defer zlibWriter.Close()

		c.responseWriter.Header().Set("Content-Encoding", "deflate")
		c.responseWriter.WriteHeader(c.statusCode)
		_, err = zlibWriter.Write(c.buffer)
		return err
	case brotliType:
		brotliWriter := brotli.NewWriterLevel(c.responseWriter, c.config.BrotliLevel)
		defer brotliWriter.Close()

		c.responseWriter.Header().Set("Content-Encoding", "br")
		c.responseWriter.WriteHeader(c.statusCode)
		_, err := brotliWriter.Write(c.buffer)
		return err
	case lzwType:
		panic("unsupported compression: LZW")
	}

	//no compression
	c.responseWriter.WriteHeader(c.statusCode)
	_, err := c.responseWriter.Write(c.buffer)
	return err
}

func (c *netHTTPCompressor) Header() http.Header {
	return c.responseWriter.Header()
}

func (c *netHTTPCompressor) Write(a []byte) (int, error) {
	c.buffer = append(c.buffer, a...)
	return len(a), nil
}

func (c *netHTTPCompressor) WriteHeader(statusCode int) {
	c.statusCode = statusCode
}

func wrapNetHTTPResponse(w http.ResponseWriter, c Config, acceptEncoding string) *netHTTPCompressor {
	return &netHTTPCompressor{
		responseWriter: w,
		config:         c,
		acceptEncoding: []byte(acceptEncoding),
		statusCode:     http.StatusOK,
	}
}

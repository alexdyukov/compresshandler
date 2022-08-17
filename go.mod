module github.com/alexdyukov/compresshandler

go 1.18

require (
	github.com/andybalholm/brotli v1.0.4
	github.com/klauspost/compress v1.15.12
	github.com/stretchr/testify v1.8.0
	github.com/valyala/fasthttp v1.43.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/alexdyukov/compresshandler/internal/encoding => ./internal/encoding

replace github.com/alexdyukov/compresshandler/internal/compressor => ./internal/compressor

replace github.com/alexdyukov/compresshandler/internal/decompressor => ./internal/decompressor

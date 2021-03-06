module github.com/alexdyukov/compresshandler

go 1.18

require (
	github.com/andybalholm/brotli v1.0.4
	github.com/klauspost/compress v1.15.6
	github.com/stretchr/testify v1.7.2
)

require (
	github.com/davecgh/go-spew v1.1.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/alexdyukov/compresshandler/internal/encoding => ./internal/encoding

replace github.com/alexdyukov/compresshandler/internal/compressors => ./internal/compressors

replace github.com/alexdyukov/compresshandler/internal/decompressors => ./internal/decompressors

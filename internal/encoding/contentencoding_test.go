package encoding_test

import (
	"testing"

	"github.com/alexdyukov/compresshandler/internal/encoding"
	"github.com/stretchr/testify/assert"
)

func TestContentEncoding(t *testing.T) {
	tests := []struct {
		name            string
		contentEncoding string
		encodingChain   []int
	}{
		{
			name:            "empty",
			contentEncoding: "",
			encodingChain:   []int{},
		},
		{
			name:            "gzip only",
			contentEncoding: "gzip",
			encodingChain:   []int{encoding.GzipType},
		},
		{
			name:            "gzip followed by deflate",
			contentEncoding: "gzip, deflate",
			encodingChain:   []int{encoding.GzipType, encoding.DeflateType},
		},
		{
			name:            "gzip followed by br",
			contentEncoding: "gzip, br",
			encodingChain:   []int{encoding.GzipType, encoding.BrType},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			contentEncoding := []byte(test.contentEncoding)

			detectedTypes, err := encoding.ParseContentEncoding(contentEncoding)
			assert.Nil(t, err)

			assert.Equal(t, len(detectedTypes), len(test.encodingChain))

			for i := 0; i < len(detectedTypes); i++ {
				assert.Equal(t, test.encodingChain[i], detectedTypes[i])
			}
		})
	}
}

package compresshandler

import (
	"testing"

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
			encodingChain:   []int{gzipType},
		},
		{
			name:            "gzip followed by deflate",
			contentEncoding: "gzip, deflate",
			encodingChain:   []int{gzipType, deflateType},
		},
		{
			name:            "gzip followed by br",
			contentEncoding: "gzip, br",
			encodingChain:   []int{gzipType, brType},
		},
	}

	// run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contentEncoding := []byte(tt.contentEncoding)
			pos := 0
			detectedType := 0

			for i := 0; i < len(tt.encodingChain); i++ {
				detectedType, pos = getEncodingType(contentEncoding, pos)
				assert.Equal(t, tt.encodingChain[i], detectedType)
			}
		})
	}
}

func TestAcceptEncoding(t *testing.T) {
	tests := []struct {
		name             string
		acceptEncoding   string
		preferedEncoding int
	}{
		{
			name:             "no compression",
			acceptEncoding:   "",
			preferedEncoding: identityType,
		},
		{
			name:             "gzip",
			acceptEncoding:   "gzip",
			preferedEncoding: gzipType,
		},
		{
			name:             "gzip and deflate",
			acceptEncoding:   "gzip, deflate",
			preferedEncoding: gzipType,
		},
		{
			name:             "quality string with br and deflate prefer",
			acceptEncoding:   "deflate, br;q=1.0, *;q=0.5",
			preferedEncoding: brType,
		},
		{
			name:             "quality string with br prefer",
			acceptEncoding:   "br;q=1.0, gzip;q=0.8, *;q=0.1",
			preferedEncoding: brType,
		},
	}

	// run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pref := getPreferedCompression([]byte(tt.acceptEncoding))
			assert.Equal(t, pref, tt.preferedEncoding)
		})
	}
}

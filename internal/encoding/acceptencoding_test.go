package encoding_test

import (
	"testing"

	"github.com/alexdyukov/compresshandler/internal/encoding"
	"github.com/stretchr/testify/assert"
)

func TestAcceptEncoding(t *testing.T) {
	tests := []struct {
		name             string
		acceptEncoding   string
		preferedEncoding int
	}{
		{
			name:             "no compression",
			acceptEncoding:   "",
			preferedEncoding: encoding.IdentityType,
		},
		{
			name:             "gzip",
			acceptEncoding:   "gzip",
			preferedEncoding: encoding.GzipType,
		},
		{
			name:             "gzip and deflate",
			acceptEncoding:   "gzip, deflate",
			preferedEncoding: encoding.GzipType,
		},
		{
			name:             "quality string with br and deflate prefer",
			acceptEncoding:   "deflate, br;q=1.0, *;q=0.5",
			preferedEncoding: encoding.BrType,
		},
		{
			name:             "quality string with br prefer",
			acceptEncoding:   "br;q=1.0, gzip;q=0.8, *;q=0.1",
			preferedEncoding: encoding.BrType,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			pref := encoding.ParseAcceptEncoding([]byte(test.acceptEncoding))

			assert.Equal(t, pref, test.preferedEncoding)
		})
	}
}

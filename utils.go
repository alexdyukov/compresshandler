package compresshandler

import "strings"

const (
	identityType = iota // Accept-Encoding: identity
	gzipType            // Accept-Encoding: gzip
	lzwType             // Accept-Encoding: compress
	zlibType            // Accept-Encoding: deflate
	brotliType          // Accept-Encoding: br
)

var (
	gzipSkipPosition     = len("gzip") - strings.IndexByte("gzip", 'g')
	identitySkipPosition = len("identity") - strings.IndexByte("identity", 'y')
	lzwSkipPosition      = len("compress") - strings.IndexByte("compress", 'c')
	zlibSkipPosition     = len("deflate") - strings.IndexByte("deflate", 'd')
	brotliSkipPosition   = len("br") - strings.IndexByte("br", 'b')
)

func parseInt(a []byte, pos int) (ans int, newpos int) {
	for pos < len(a) && a[pos] >= '0' && a[pos] <= '9' {
		ans = 10*ans + int(a[pos]-'0')
		pos += 1
	}
	return ans, pos
}

func parseFloat(a []byte, pos int) (ans float64, newpos int) {
	// 000000000.00000000
	// left		.	right
	var leftInt, rightInt int
	var leftFloat, rightFloat float64

	leftInt, pos = parseInt(a, pos)
	leftFloat = float64(leftInt)

	if pos >= len(a) || a[pos] != '.' {
		return leftFloat, pos
	}

	rightInt, pos = parseInt(a, pos+1)
	rightFloat = float64(rightInt)
	for rightFloat >= 1.0 {
		rightFloat /= 10.0
	}

	return leftFloat + rightFloat, pos
}

func abs(v float64) float64 {
	if v < 0.0 {
		return -v
	}
	return v
}

func parseEncoding(a []byte, pos int) (t int, q float64, newpos int) {
	q = -1.0
	t = identityType

	//https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Accept-Encoding
	for pos < len(a) && (a[pos] != ',') {
		switch a[pos] {
		case 'g', 'G': // g stands only in gzip
			t = gzipType
			pos += gzipSkipPosition
		case 'c', 'C': // c stands only in compress
			t = lzwType
			pos += lzwSkipPosition
		case 'f', 'F': // f stands only in deflate
			t = zlibType
			pos += zlibSkipPosition
		case 'b', 'B': // b stands only in br
			t = brotliType
			pos += brotliSkipPosition
		case 'y', 'Y': // y stands only in identity
			t = identityType
			pos += identitySkipPosition
		case '*': // it means any other, but we cast gzip as web standart
			t = gzipType
			pos += 1
		case '0', '1': // possible values in range [0.0;1.0] thats why we are looking for 0 and 1 only
			q, pos = parseFloat(a, pos)
		default:
			pos += 1
		}
	}

	return t, abs(q), pos + 1
}

func getPreferedCompression(acceptEncoding []byte) int {
	bestType := identityType
	bestQuality := 0.0
	t, q, pos := 0, 0.0, 0

	for pos < len(acceptEncoding) {
		t, q, pos = parseEncoding(acceptEncoding, pos)
		// TODO LZW; unsupported now
		if t == lzwType {
			continue
		}
		if bestQuality < q {
			bestType, bestQuality = t, q
		}
	}
	return bestType
}

func getRequestCompression(contentEncoding []byte) []int {
	contentEncodingTypes := make([]int, 0, 4)

	for pos := 0; pos < len(contentEncoding); pos += 1 {
		switch contentEncoding[pos] {
		case 'z', 'Z': // z stands only in gzip >> gzip
			contentEncodingTypes = append(contentEncodingTypes, gzipType)
		case 'c', 'C': // c stands only in compress >> lzw
			contentEncodingTypes = append(contentEncodingTypes, lzwType)
		case 'f', 'F': // f stands only in deflate >> zlib
			contentEncodingTypes = append(contentEncodingTypes, zlibType)
		case 'b', 'B': // b stands only in br >> brotli
			contentEncodingTypes = append(contentEncodingTypes, brotliType)
		}
	}

	return contentEncodingTypes
}

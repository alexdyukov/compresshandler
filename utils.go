package compresshandler

const (
	identityType = iota // Accept-Encoding: identity
	gzipType            // Accept-Encoding: gzip
	lzwType             // Accept-Encoding: compress
	zlibType            // Accept-Encoding: deflate
	brotliType          // Accept-Encoding: br
)

func parseInt(a []byte, pos int) (int, int) {
	ans := 0
	for pos < len(a) && a[pos] >= '0' && a[pos] <= '9' {
		ans = 10*ans + int(a[pos]-'0')
		pos++
	}

	return ans, pos
}

func parseFloat(str []byte, pos int) (float64, int) {
	// float example: 000000.000000
	// how's names go: left . right
	leftInt, rightInt := 0, 0
	leftFloat, rightFloat := 0.0, 0.0

	leftInt, pos = parseInt(str, pos)
	leftFloat = float64(leftInt)

	if pos >= len(str) || str[pos] != '.' {
		return leftFloat, pos
	}

	rightInt, pos = parseInt(str, pos+1)
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

func parseEncoding(str []byte, pos int) (int, float64, int) {
	encodingType, encodingQuality := identityType, -1.0

	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Accept-Encoding
	for pos < len(str) && (str[pos] != ',') {
		switch str[pos] {
		case 'g', 'G': // g stands only in gzip
			encodingType = gzipType
		case 'c', 'C': // c stands only in compress
			encodingType = lzwType
		case 'f', 'F': // f stands only in deflate
			encodingType = zlibType
		case 'b', 'B': // b stands only in br
			encodingType = brotliType
		case 'y', 'Y': // y stands only in identity
			encodingType = identityType
		case '*': // it means any other, but we cast gzip as web standart
			encodingType = gzipType
		case '0', '1': // possible values in range [0.0;1.0] thats why we are looking for 0 and 1 only
			encodingQuality, pos = parseFloat(str, pos)
		}
		pos++
	}

	return encodingType, abs(encodingQuality), pos + 1
}

func getPreferedCompression(acceptEncoding []byte) int {
	bestType := identityType
	bestQuality := 0.0
	currentType, currentQuality, currentPos := 0, 0.0, 0

	for currentPos < len(acceptEncoding) {
		currentType, currentQuality, currentPos = parseEncoding(acceptEncoding, currentPos)

		if currentType == lzwType {
			continue
		}

		if bestQuality < currentQuality {
			bestType, bestQuality = currentType, currentQuality
		}
	}

	return bestType
}

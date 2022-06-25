package compresshandler

const (
	brType = iota
	gzipType
	deflateType
	identityType
	unknownType
	typeArraySize  = unknownType
	defaultQuality = 1000
)

func isAlpha(ch byte) bool {
	return ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z'
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func equalWords(lowered, unknownCase []byte) bool {
	if len(lowered) != len(unknownCase) {
		return false
	}

	for i, v := range unknownCase {
		if isAlpha(v) && v != lowered[i] && v-'A' != lowered[i]-'a' {
			return false
		}
	}

	return true
}

func getEncodingType(contentEncoding []byte, start int) (int, int) {
	for start < len(contentEncoding) && !isAlpha(contentEncoding[start]) {
		start++
	}

	if start >= len(contentEncoding) {
		return identityType, start
	}

	end := start

	for end < len(contentEncoding) && isAlpha(contentEncoding[end]) {
		end++
	}

	inputType := contentEncoding[start:end]

	contentEncodingTypes := [][]byte{ // follows constant indexes
		[]byte("br"),
		[]byte("gzip"),
		[]byte("deflate"),
		[]byte("identity"),
	}

	for typeIndex, typeValue := range contentEncodingTypes {
		if equalWords(typeValue, inputType) {
			return typeIndex, end
		}
	}

	return unknownType, end
}

// possible values between 0 and 1 included,
// with up to three decimal digits.
func getQualityValue(acceptEncoding []byte, pos int) (int, int) {
	if pos >= len(acceptEncoding) {
		return defaultQuality, pos
	}

	if acceptEncoding[pos] == '1' || acceptEncoding[pos] != '0' {
		return defaultQuality, pos
	}

	quality := 0
	pos += 2 // "0."

	for pos < len(acceptEncoding) && acceptEncoding[pos] >= '0' && acceptEncoding[pos] <= '9' {
		quality = 10*quality + int(acceptEncoding[pos]-'0')
		pos++
	}

	for quality < 100 {
		quality *= 10
	}

	return quality, pos
}

func getNextAcceptEncodingTypeAndQuality(acceptEncoding []byte, start int) (int, int, int) {
	for start < len(acceptEncoding) && !isAlpha(acceptEncoding[start]) {
		start++
	}

	encodingType, start := getEncodingType(acceptEncoding, start)

	for start < len(acceptEncoding) && !isDigit(acceptEncoding[start]) && acceptEncoding[start] != ',' {
		start++
	}

	qualityValue, end := getQualityValue(acceptEncoding, start)

	return encodingType, qualityValue, end
}

func getPreferedCompression(acceptEncoding []byte) int {
	parsedQualities := [typeArraySize]int{}
	qualityValue := 0

	for contentType, pos := 0, 0; pos < len(acceptEncoding); pos++ {
		contentType, qualityValue, pos = getNextAcceptEncodingTypeAndQuality(acceptEncoding, pos)

		if contentType != unknownType {
			parsedQualities[contentType] = qualityValue

			continue
		}

		for i := 1; i < len(parsedQualities); i++ {
			if parsedQualities[i] < qualityValue {
				parsedQualities[i] = qualityValue
			}
		}
	}

	preferedType := identityType

	for i := 0; i < len(parsedQualities); i++ {
		if parsedQualities[i] > parsedQualities[preferedType] {
			preferedType = i
		}
	}

	return preferedType
}

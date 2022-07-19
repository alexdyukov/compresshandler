package encoding

const (
	BrType = iota
	GzipType
	DeflateType
	IdentityType
	UnknownType
	typeArraySize  = UnknownType
	defaultQuality = 1000
)

func isAlpha(ch byte) bool {
	return ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z'
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func parseQuality(str []byte, pos int) (int, int) {
	quality := 0
	for i := 0; i < 3; i++ {
		quality *= 10
		if pos < len(str) && isDigit(str[pos]) {
			quality += int(str[pos] - '0')
			pos++
		}
	}

	return quality, pos
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

func getNextEncodingType(header []byte, start int) (int, int) {
	for start < len(header) && !isAlpha(header[start]) {
		start++
	}

	if start >= len(header) {
		return IdentityType, start
	}

	end := start

	for end < len(header) && isAlpha(header[end]) {
		end++
	}

	inputType := header[start:end]

	contentEncodingTypes := [][]byte{
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

	return UnknownType, end
}

// possible values between 0 and 1 included,
// with up to three decimal digits.
func getNextQualityValue(header []byte, pos int) (int, int) {
	for pos < len(header) && !isDigit(header[pos]) && header[pos] != ',' {
		pos++
	}

	if pos >= len(header) {
		return defaultQuality, pos
	}

	if header[pos] == '1' || header[pos] != '0' {
		return defaultQuality, pos
	}

	pos += 2

	return parseQuality(header, pos)
}

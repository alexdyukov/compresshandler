package encoding

// ParseAcceptEncoding parses Accept-Encoding http header value.
func ParseAcceptEncoding(headerValue []byte) int {
	parsedQualities := [typeArraySize]int{}

	var (
		encodingType int
		qualityValue int
	)

	for pos := 0; pos < len(headerValue); pos++ {
		encodingType, pos = getNextEncodingType(headerValue, pos)
		qualityValue, pos = getNextQualityValue(headerValue, pos)

		if encodingType != UnknownType {
			parsedQualities[encodingType] = qualityValue

			continue
		}

		for i := 1; i < len(parsedQualities); i++ {
			if parsedQualities[i] < qualityValue {
				parsedQualities[i] = qualityValue
			}
		}
	}

	preferedType := IdentityType

	for i := 0; i < len(parsedQualities); i++ {
		if parsedQualities[i] > parsedQualities[preferedType] {
			preferedType = i
		}
	}

	return preferedType
}

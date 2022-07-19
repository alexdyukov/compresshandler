package encoding

func ParseAcceptEncoding(header []byte) int {
	parsedQualities := [typeArraySize]int{}
	encodingType := 0
	qualityValue := 0

	for pos := 0; pos < len(header); pos++ {
		encodingType, pos = getNextEncodingType(header, pos)
		qualityValue, pos = getNextQualityValue(header, pos)

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

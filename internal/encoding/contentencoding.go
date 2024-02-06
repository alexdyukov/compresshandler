package encoding

// UnsupportedContentEncodingError indicates that a provided Content-Encoding field is not supported.
type UnsupportedContentEncodingError string

func (err UnsupportedContentEncodingError) Error() string {
	return "unsupported content encoding: " + string(err)
}

// ParseContentEncoding parses Content-Encoding http header value.
func ParseContentEncoding(headerValue []byte) ([]int, error) {
	var (
		contentType int
		encodings   = []int{}
	)

	for pos := 0; pos < len(headerValue); pos++ {
		contentType, pos = getNextEncodingType(headerValue, pos)
		if contentType == UnknownType {
			return encodings, UnsupportedContentEncodingError(headerValue)
		}

		encodings = append(encodings, contentType)
	}

	return encodings, nil
}

// ToString translates enum type into string value.
// Its faster to have enum, rather then redefine types with String() method.
func ToString(encoding int) string {
	switch encoding {
	case BrType:
		return "br"
	case GzipType:
		return "gzip"
	case DeflateType:
		return "deflate"
	default:
		return ""
	}
}

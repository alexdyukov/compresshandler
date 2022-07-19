package encoding

type UnsupportedContentEncodingError string

func (err UnsupportedContentEncodingError) Error() string {
	return "unsupported content encoding: " + string(err)
}

func ParseContentEncoding(header []byte) ([]int, error) {
	var (
		contentType int
		encodings   = []int{}
	)

	for pos := 0; pos < len(header); pos++ {
		contentType, pos = getNextEncodingType(header, pos)
		if contentType == UnknownType {
			return encodings, UnsupportedContentEncodingError(header)
		}

		encodings = append(encodings, contentType)
	}

	return encodings, nil
}

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

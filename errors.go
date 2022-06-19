package compresshandler

type CompressError struct {
	ErrorString string
}

func (ce *CompressError) Error() string {
	return ce.ErrorString
}

var (
	ErrNotSupported = &CompressError{"Unsupported content encoding"}
)

package constants

import "errors"

var (
	// ErrorInvalidFlow returned while flow is invalid.
	ErrorInvalidFlow = errors.New("input source and destination is invalid")
	// ErrorFileTooLarge returned while file is too large.
	ErrorFileTooLarge = errors.New("file is too large")
)

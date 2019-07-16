package constants

import "errors"

var (
	// ErrorFlowInvalid returned while flow is invalid.
	ErrorFlowInvalid = errors.New("flow is invalid")
	// ErrorActionNotImplemented returned while current is not implemented.
	ErrorActionNotImplemented = errors.New("action not implemented")

	// ErrorExpectSizeRequired returned while expect size is required but not given.
	ErrorExpectSizeRequired = errors.New("expect-size is required")
	// ErrorByteSizeInvalid returned while byte size is invalid.
	ErrorByteSizeInvalid = errors.New("byte size is invalid")
	// ErrorZoneRequired returned while zone is required but not given.
	ErrorZoneRequired = errors.New("zone is required")

	// ErrorExternalServiceError returned while external service error happened.
	ErrorExternalServiceError = errors.New("external service error happened")

	// ErrorQsPathAccessForbidden returned while qingstor path access is forbidden.
	ErrorQsPathAccessForbidden = errors.New("qingstor path access forbidden")
	// ErrorQsPathNotFound returned while qingstor path is not found.
	ErrorQsPathNotFound = errors.New("qingstor path not found")
	// ErrorQsPathInvalid returned while qs-path is invalid.
	ErrorQsPathInvalid = errors.New("qingstor path invalid")
	// ErrorQsPathObjectKeyRequired returned while object key is required but not given.
	ErrorQsPathObjectKeyRequired = errors.New("qingstor path object key is required")

	// ErrorFileTooLarge returned while file is too large.
	ErrorFileTooLarge = errors.New("file too large")
	// ErrorFileNotExist returned while file is not found.
	ErrorFileNotExist = errors.New("file not exist")
	// ErrorFileSizeInvalid returned while get file size not valid.
	ErrorFileSizeInvalid = errors.New("file size invalid")
)

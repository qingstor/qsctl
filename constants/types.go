package constants

// DirectoryContentType is the content type qingstor used for directory.
const DirectoryContentType = "application/x-directory"

const (
	// MaximumObjectSize is the max object size for a single object, 50TB.
	MaximumObjectSize = 50 * 1024 * 1024 * 1024 * 1024
	// DefaultMultipartBoundarySize is the default boundary size for multipart, 1GB.
	DefaultMultipartBoundarySize = 1024 * 1024 * 1024
)

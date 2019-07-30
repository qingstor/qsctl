package constants

// DirectoryContentType is the content type qingstor used for directory.
const DirectoryContentType = "application/x-directory"

const (
	// ACLDirectory is the acl of qingstor x-directory
	ACLDirectory = "drwxr-xr-x"
	// ACLObject is the acl of qingstor object
	ACLObject = "-rwxr-xr-x"
)

const (
	// MaximumMultipartNumber is the max multipart number.
	MaximumMultipartNumber = 10000
	// MaximumPartSize is the max part size for single part, 5GB.
	MaximumPartSize = 5 * 1024 * 1024 * 1024
	// MaximumObjectSize is the max object size for a single object, 50TB.
	MaximumObjectSize = MaximumMultipartNumber * MaximumPartSize
	// MaximumAutoMultipartSize is the size limit for auto part size detect.
	MaximumAutoMultipartSize = MaximumPartSize / 5
	// DefaultPartSize is the default part size, 128MB.
	DefaultPartSize = 128 * 1024 * 1024
)

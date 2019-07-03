package action

import (
	"io"
	"os"
	"strings"

	"github.com/yunify/qsctl/constants"
	"github.com/yunify/qsctl/contexts"
)

// ParseDirection will parse the data direction
func ParseDirection(src, dst string) (flow string, err error) {
	// If src and dst both local file or both remote object, the path is invalid.
	if strings.HasPrefix(src, "qs://") == strings.HasPrefix(dst, "qs://") {
		return "", constants.ErrorInvalidFlow
	}

	if strings.HasPrefix(src, "qs://") {
		return constants.DirectionRemoteToLocal, nil
	}
	return constants.DirectionLocalToRemote, nil
}

// ParseFilePathForRead will parse file path and open an io.Reader for read.
func ParseFilePathForRead(filePath string) (r io.Reader, err error) {
	// Use - means we will read from stdin.
	if filePath == "-" {
		return os.Stdin, nil
	}

	_, err = os.Stat(filePath)
	if err != nil {
		panic(err)
	}

	return os.Open(filePath)
}

// ParseFilePathForWrite will parse a file path and open an io.Write for write.
func ParseFilePathForWrite() {}

// ParseQsPathForRead will parse a qs path and open an io.Reader for read.
func ParseQsPathForRead() {}

// ParseQsPathForWrite will parse a qs path for write.
func ParseQsPathForWrite(remotePath string) (objectKey string, err error) {
	// qs://abc/xyz -> []string{"qs:", "", "abc", "xyz"}
	p := strings.Split(remotePath, "/")
	if p[0] != "qs:" || p[1] != "" || p[2] == "" {
		// FIXME
		panic("invalid qs path")
	}
	bucketName := p[2]

	_, err = contexts.SetupBuckets(bucketName, "")
	if err != nil {
		// FIXME
		panic(err)
	}

	// FIXME: user may input qs://abc
	// Trim "qs://" + bucketName + "/"
	objectKey = remotePath[5+len(bucketName)+1:]
	return
}

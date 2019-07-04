package action

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/yunify/qsctl/constants"
	"github.com/yunify/qsctl/contexts"
)

// ParseDirection will parse the data direction
func ParseDirection(src, dst string) (flow string, err error) {
	// If src and dst both local file or both remote object, the path is invalid.
	if strings.HasPrefix(src, "qs://") == strings.HasPrefix(dst, "qs://") {
		log.Errorf("Action between <%s> and <%s> is invalid", src, dst)
		return "", constants.ErrorFlowInvalid
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
	if os.IsNotExist(err) {
		log.Infof("File <%s> is not exist, please check your input path", filePath)
		return nil, constants.ErrorFileNotExist
	}
	if err != nil {
		log.Errorf("action: Stat file failed [%s]", err)
		return
	}

	return os.Open(filePath)
}

// ParseFilePathForWrite will parse a file path and open an io.Write for write.
func ParseFilePathForWrite(filePath string) (w io.Writer, err error) {
	// Use - means we will read from stdin.
	if filePath == "-" {
		return os.Stdout, nil
	}

	// Create dir automatically.
	err = os.MkdirAll(filepath.Dir(filePath), os.ModeDir|0664)
	if err != nil {
		log.Errorf("Mkdir <%s> failed [%v]", filePath, err)
		return nil, err
	}

	return os.Create(filePath)
}

// ParseQsPath will parse a qs path and prepare a bucket.
func ParseQsPath(remotePath string, objectKeyRequired bool) (objectKey string, err error) {
	// qs://abc/xyz -> []string{"qs:", "", "abc", "xyz"}
	p := strings.Split(remotePath, "/")
	if p[0] != "qs:" || p[1] != "" || p[2] == "" {
		log.Infof("<%s> is not a valid qingstor path", remotePath)
		return "", constants.ErrorQsPathInvalid
	}
	bucketName := p[2]

	_, err = contexts.SetupBuckets(bucketName, "")
	if err != nil {
		return
	}

	if len(p) >= 4 {
		// Trim "qs://" + bucketName + "/"
		objectKey = remotePath[5+len(bucketName)+1:]
		return
	}

	if objectKeyRequired {
		return "", constants.ErrorQsPathObjectKeyRequired
	}
	// Handle user input "qs://abc"
	return "", nil
}

// CalculateConcurrentWorkers will calculate the current workers via limit and part size.
func CalculateConcurrentWorkers(partSize int64) (n int) {
	// If the part size is over the limit, we will only use one worker.
	if contexts.MaximumMemoryContent <= partSize {
		return 1
	}

	return int(contexts.MaximumMemoryContent / partSize)
}

// CalculatePartSize will calculate the object's part size.
func CalculatePartSize(size int64) (partSize int64, err error) {
	partSize = constants.DefaultPartSize

	if size > constants.MaximumObjectSize {
		log.Errorf("File with size <%d> is too large", size)
		return 0, constants.ErrorFileTooLarge
	}

	for size/partSize >= int64(constants.MaximumMultipartNumber) {
		if partSize < constants.MaximumAutoMultipartSize {
			partSize = partSize << 1
			continue
		}
		// Try to adjust partSize if it is too small and account for
		// integer division truncation.
		partSize = size/int64(constants.MaximumMultipartNumber) + 1
		break
	}
	return
}

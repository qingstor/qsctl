package utils

import (
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/task/types"
)

// ParseFilePathForRead will parse file path and open an io.Reader for read.
func ParseFilePathForRead(filePath string) (r io.Reader, err error) {
	// Use - means we will read from stdin.
	if filePath == "-" {
		return os.Stdin, nil
	}

	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		log.Infof("File <%s> is not exist, please check your input", filePath)
		return nil, constants.ErrorFileNotExist
	}
	if err != nil {
		log.Errorf("Stat file failed [%s]", err)
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

// ParseQsPath will parse a qs path.
func ParseQsPath(remotePath string) (bucketName, objectKey string, err error) {
	// qs-path includes three part: "qs://" prefix, bucket name and object key.
	// "qs://" prefix could be emit.
	pattern := "^(?:qs://)?([a-z\\d][a-z-\\d]{4,61}[a-z\\d])?(.*)?$"

	x := regexp.MustCompile(pattern).FindStringSubmatch(remotePath)
	if len(x) != 3 || x[1] == "" {
		return "", "", constants.ErrorQsPathInvalid
	}

	bucketName, objectKey = x[1], x[2]

	// TODO: add bucket name and object key check here.

	// Trim all left "/"
	objectKey = strings.TrimLeft(objectKey, "/")

	return bucketName, objectKey, nil
}

// CalculateConcurrentWorkers will calculate the current workers via limit and part size.
func CalculateConcurrentWorkers(partSize, maximumMemory int64) (n int) {
	// If the part size is over the limit, we will only use one worker.
	if maximumMemory <= partSize {
		return 1
	}

	return int(maximumMemory / partSize)
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

// CalculateSeekableFileSize will calculate the seekable file's size.
func CalculateSeekableFileSize(r io.Seeker) (size int64, err error) {
	// Move the start to make sure size read correctly.
	_, err = r.Seek(0, io.SeekStart)
	if err != nil {
		return
	}

	size, err = r.Seek(0, io.SeekEnd)
	if err != nil {
		return
	}
	return
}

// SubmitNextTask will fetch next todo and submit to pool.
func SubmitNextTask(t interface {
	types.Todoist
	types.PoolGetter
}) {
	fn := t.NextTODO()
	if fn == nil {
		return
	}

	t.GetPool().Submit(fn(t))
}

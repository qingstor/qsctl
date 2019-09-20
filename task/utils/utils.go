package utils

import (
	"io"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/task/types"
)

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

// CalculateFileSize will calculate the seekable file's size.
func CalculateFileSize(r io.Seeker) (size int64, err error) {
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

// SubmitNextTask will fetch next todo and submit to pool async.
func SubmitNextTask(t types.Tasker) {
	fn := t.NextTODO()
	if fn == nil {
		return
	}

	pool := t.GetPool()
	if pool.Free() > 0 {
		pool.Submit(fn(t))
		return
	}
	go pool.Submit(fn(t))
}

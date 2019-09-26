package utils

import (
	"fmt"
	"io"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/pkg/fault"
	"github.com/yunify/qsctl/v2/pkg/types"
)

// CalculatePartSize will calculate the object's part size.
func CalculatePartSize(size int64) (partSize int64, err error) {
	partSize = constants.DefaultPartSize

	if size > constants.MaximumObjectSize {
		return 0, fmt.Errorf("calculate part size failed: {%w}", fault.NewLocalFileTooLarge(nil, size))
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
	// We will update task's fault in run(), and check it every time we submit it to the pool.
	if t.ValidateFault() {
		return
	}

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

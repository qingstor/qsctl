package task

import (
	"bufio"
	"bytes"
	"io"
	"sync"
	"sync/atomic"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/pkg/fault"
	"github.com/yunify/qsctl/v2/pkg/types"
	"github.com/yunify/qsctl/v2/task/common"
)

// NewCopyStreamTask will create a copy stream task.
func (t *CopyStreamTask) new() {
	bytesPool := &sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 0, t.GetPartSize()))
		},
	}
	t.SetBytesPool(bytesPool)

	// TODO: we will use expect size to calculate part size later.
	t.SetPartSize(constants.DefaultPartSize)

	t.SetScheduler(types.NewScheduler(NewCopyPartialStreamTask))

	currentOffset := int64(0)
	t.SetCurrentOffset(&currentOffset)

	// We don't know how many data in stream, set it to -1 as an indicate.
	// We will set current offset to -1 when got an EOF from stream.
	t.SetTotalSize(-1)

	t.AddTODOs(
		common.NewMultipartInitTask,
		common.NewWaitTask,
		common.NewMultipartCompleteTask,
	)
	return
}

// NewCopyPartialStreamTask will create a new Task.
func (t *CopyPartialStreamTask) new() {
	// Set size and update offset.
	partSize := t.GetPartSize()
	r := bufio.NewReader(io.LimitReader(t.GetStream(), partSize))
	b := t.GetBytesPool().Get().(*bytes.Buffer)
	n, err := io.Copy(b, r)
	if err != nil {
		t.TriggerFault(fault.NewUnhandled(err))
		return
	}

	t.SetOffset(*t.GetCurrentOffset())
	t.SetSize(n)
	t.SetContent(b)
	if n < partSize {
		// Set current offset to -1 to mark the stream has been drain out.
		atomic.StoreInt64(t.GetCurrentOffset(), -1)
	} else {
		atomic.AddInt64(t.GetCurrentOffset(), t.GetSize())
	}

	t.AddTODOs(
		common.NewStreamMD5SumTask,
		common.NewMultipartStreamUploadTask,
	)
	return
}

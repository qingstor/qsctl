package task

import (
	"bufio"
	"bytes"
	"io"
	"sync"
	"sync/atomic"

	"github.com/Xuanwo/navvy"
	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/task/common"
	"github.com/yunify/qsctl/v2/task/types"
)

// NewCopyStreamTask will create a copy stream task.
func NewCopyStreamTask(task types.Todoist) navvy.Task {
	t, _ := initCopyStreamTask(task)

	bytesPool := &sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 0, t.GetPartSize()))
		},
	}
	t.SetBytesPool(bytesPool)

	// TODO: we will use expect size to calculate part size later.
	t.SetPartSize(constants.DefaultPartSize)

	t.SetWaitGroup(&sync.WaitGroup{})

	t.SetTaskConstructor(NewCopyPartialStreamTask)

	currentPartNumber := int32(0)
	t.SetCurrentPartNumber(&currentPartNumber)

	currentOffset := int64(0)
	t.SetCurrentOffset(&currentOffset)

	// We don't know how many data in stream, set it to -1 as an indicate.
	// We will set current offset to -1 when got an EOF from stream.
	t.SetSize(-1)

	t.AddTODOs(
		common.NewMultipartInitTask,
		common.NewWaitTask,
		common.NewMultipartCompleteTask,
	)
	return t
}

// NewCopyPartialStreamTask will create a new Task.
func NewCopyPartialStreamTask(task types.Todoist) navvy.Task {
	t, o := initCopyPartialStreamTask(task)

	// Set part number and update current part number.
	currentPartNumber := o.GetCurrentPartNumber()
	t.SetPartNumber(int(*currentPartNumber))
	atomic.AddInt32(currentPartNumber, 1)

	// Set size and update offset.
	partSize := o.GetPartSize()
	r := bufio.NewReader(io.LimitReader(o.GetStream(), partSize))
	b := o.GetBytesPool().Get().(*bytes.Buffer)
	n, err := io.Copy(b, r)
	if err != nil {
		panic(err)
	}

	t.SetSize(n)
	t.SetContent(b)
	if n < partSize {
		// Set current offset to -1 to mark the stream has been drain out.
		atomic.StoreInt64(o.GetCurrentOffset(), -1)
	} else {
		atomic.AddInt64(o.GetCurrentOffset(), t.GetSize())
	}

	t.AddTODOs(
		common.NewStreamMD5SumTask,
		common.NewMultipartStreamUploadTask,
	)
	return t
}

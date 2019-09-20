package task

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"sync"
	"sync/atomic"

	"github.com/Xuanwo/navvy"
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

	t.SetWaitGroup(&sync.WaitGroup{})

	var err error

	f := os.Stdin
	if t.GetPath() != "-" {
		f, err = os.Open(t.GetPath())
		if err != nil {
			panic(err)
		}
	}
	t.SetStream(f)
	t.SetTaskConstructor(NewCopyPartialStreamTask)

	currentPartNumber := int32(0)
	t.SetCurrentPartNumber(&currentPartNumber)

	currentOffset := int64(0)
	t.SetCurrentOffset(&currentOffset)

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
	atomic.AddInt64(o.GetCurrentOffset(), t.GetSize())

	t.AddTODOs(
		common.NewStreamMD5SumTask,
		common.NewMultipartStreamUploadTask,
	)
	return t
}

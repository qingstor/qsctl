package task

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"sync"
	"sync/atomic"

	"github.com/Xuanwo/navvy"
	log "github.com/sirupsen/logrus"
	"github.com/yunify/qsctl/v2/storage"
	"github.com/yunify/qsctl/v2/task/common"
	"github.com/yunify/qsctl/v2/task/types"
	"github.com/yunify/qsctl/v2/task/utils"
)

// CopyStreamTask will copy a stream to remote.
type CopyStreamTask struct {
	// Input value
	types.FilePath
	types.ObjectKey
	types.Storage

	// Runtime value
	types.Todo
	types.Pool
	types.Size
	types.UploadID
	types.WaitGroup
	types.CurrentPartNumber
	types.CurrentOffset
	types.PartSize
	types.Stream
	types.BytesPool
	types.TaskConstructor
}

type CopyStreamTaskOptions struct {
}

// Run implement navvy.Task
func (t *CopyStreamTask) Run() {
	log.Debugf("Task <%s> for Object <%s> started.", "CopyStreamTask", t.GetObjectKey())

	utils.SubmitNextTask(t)
}

func NewCopyStreamTask(objectKey string, storage storage.ObjectStorage) *CopyStreamTask {
	t := &CopyStreamTask{}
	t.SetObjectKey(objectKey)
	t.SetStorage(storage)

	// TODO: decide part size

	bytesPool := &sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 0, t.GetPartSize()))
		},
	}
	t.SetBytesPool(bytesPool)

	var err error

	f := os.Stdin
	if t.GetFilePath() != "-" {
		f, err = os.Open(t.GetFilePath())
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

// CopyPartialStreamTask will execute CopyPartialStream Task
type CopyPartialStreamTask struct {
	types.ObjectKey
	types.UploadID
	types.PartNumber
	types.Size
	types.Offset
	types.Storage

	types.Todo
	types.Content
	types.MD5Sum
	types.WaitGroup
	types.Pool
}

// Run implement navvy.Task
func (t *CopyPartialStreamTask) Run() {
	utils.SubmitNextTask(t)
}

// NewCopyPartialStreamTask will create a new Task.
func NewCopyPartialStreamTask(task types.Todoist) navvy.Task {
	o, ok := task.(*CopyStreamTask)
	if !ok {
		panic("parent task is not a CopyStreamTask")
	}

	t := &CopyPartialStreamTask{}
	t.SetUploadID(o.GetUploadID())
	t.SetWaitGroup(o.GetWaitGroup())
	t.SetObjectKey(o.GetObjectKey())

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

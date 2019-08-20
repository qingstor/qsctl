package task

import (
	"github.com/Xuanwo/navvy"
	log "github.com/sirupsen/logrus"

	"github.com/yunify/qsctl/v2/storage"
	"github.com/yunify/qsctl/v2/task/common"
	"github.com/yunify/qsctl/v2/task/types"
	"github.com/yunify/qsctl/v2/task/utils"
)

// CopyFileTask will handle all copy file task
type CopyFileTask struct {
	// Input value
	types.FilePath
	types.ObjectKey
	types.Storage

	// Runtime value
	types.Todo
	types.Pool
}

// NewCopyFileTask will create a new copy file task.
func NewCopyFileTask(filePath, objectKey string, storage storage.ObjectStorage) *CopyFileTask {
	t := &CopyFileTask{}
	t.SetFilePath(filePath)
	t.SetObjectKey(objectKey)

	pool, err := navvy.NewPool(10)
	if err != nil {
		panic(err)
	}
	t.SetPool(pool)

	t.SetStorage(storage)
	return t
}

// Run implement navvy.Task
func (t *CopyFileTask) Run() {
	// TODO: should change file's size to use different task.
	t.GetPool().Submit(NewCopyLargeFileTask(t))
}

// CopyStreamTask will copy a stream to remote.
type CopyStreamTask struct {
}

// CopyLargeFileTask will execute CopyLargeFile Task
type CopyLargeFileTask struct {
	// Input value
	types.FilePath
	types.ObjectKey
	types.Pool
	types.Storage

	// Runtime value
	types.Todo
	types.UploadID
	types.TotalParts
	types.WaitGroup
}

// NewCopyLargeFileTask will create a new Task.
func NewCopyLargeFileTask(task types.Todoist) navvy.Task {
	o, ok := task.(*CopyFileTask)
	if !ok {
		panic("parent task is not a CopyFileTask")
	}

	log.Debugf("Start copy large file task")

	x := &CopyLargeFileTask{}
	x.SetObjectKey(o.GetObjectKey())
	x.SetFilePath(o.GetFilePath())
	x.SetPool(o.GetPool())
	x.SetStorage(o.GetStorage())

	x.AddTODOs(
		NewMultipartObjectInitTask,
		common.NewWaitTask,
		NewMultipartObjectCompleteTask,
	)
	return x
}

// Run implement navvy.Task
func (t *CopyLargeFileTask) Run() {
	utils.SubmitNextTask(t)
}

// CopyPartialFileTask will execute CopyPartialFile Task
type CopyPartialFileTask struct {
	types.FilePath
	types.ObjectKey
	types.UploadID
	types.PartNumber
	types.Size
	types.Offset
	types.Storage

	types.Todo
	types.MD5Sum
	types.WaitGroup
	types.Pool
}

// Run implement navvy.Task
func (t *CopyPartialFileTask) Run() {
	utils.SubmitNextTask(t)
}

// NewCopyPartialFileTask will create a new Task.
func NewCopyPartialFileTask(
	objectKey, filePath, uploadID string,
	partNumber int,
	offset, size int64,
) *CopyPartialFileTask {
	t := &CopyPartialFileTask{}
	t.SetPartNumber(partNumber)
	t.SetOffset(offset)
	t.SetSize(size)
	t.SetUploadID(uploadID)
	t.SetFilePath(filePath)
	t.SetObjectKey(objectKey)

	t.AddTODOs(
		common.NewSeekableMD5SumTask,
		NewMultipartObjectUploadTask,
	)
	return t
}

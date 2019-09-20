package task

import (
	"os"
	"sync/atomic"

	"github.com/Xuanwo/navvy"
	"github.com/yunify/qsctl/v2/constants"

	"github.com/yunify/qsctl/v2/task/common"
	"github.com/yunify/qsctl/v2/task/types"
	"github.com/yunify/qsctl/v2/task/utils"
)

// NewCopyFileTask will create a new copy file task.
func NewCopyFileTask(task types.Todoist) navvy.Task {
	t, _ := initCopyFileTask(task)

	f, err := os.Open(t.GetPath())
	if err != nil {
		panic(err)
	}
	defer f.Close()

	size, err := utils.CalculateFileSize(f)
	if err != nil {
		panic(err)
	}
	t.SetSize(size)

	if size >= constants.MaximumAutoMultipartSize {
		t.AddTODOs(NewCopyLargeFileTask)
	} else {
		t.AddTODOs(NewCopySmallFileTask)
	}
	return t
}

// NewCopySmallFileTask will create a new small file task.
func NewCopySmallFileTask(task types.Todoist) navvy.Task {
	t, _ := initCopySmallFileTask(task)
	t.AddTODOs(
		common.NewFileMD5SumTask,
		common.NewFileUploadTask,
	)
	return t
}

// NewCopyLargeFileTask will create a new Task.
func NewCopyLargeFileTask(task types.Todoist) navvy.Task {
	t, _ := initCopyLargeFileTask(task)

	// Init part size.
	partSize, err := utils.CalculatePartSize(t.GetSize())
	if err != nil {
		panic(err)
	}
	t.SetPartSize(partSize)

	t.SetTaskConstructor(NewCopyPartialFileTask)

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

// NewCopyPartialFileTask will create a new Task.
func NewCopyPartialFileTask(task types.Todoist) navvy.Task {
	t, o := initCopyPartialFileTask(task)

	totalSize := o.GetSize()

	// Set part number and update current part number.
	currentPartNumber := o.GetCurrentPartNumber()
	t.SetPartNumber(int(*currentPartNumber))
	atomic.AddInt32(currentPartNumber, 1)

	// Set size and update offset.
	offset := o.GetPartSize() * int64(t.GetPartNumber())
	t.SetOffset(offset)
	if totalSize < offset {
		t.SetSize(totalSize - offset)
	} else {
		t.SetSize(o.GetPartSize())
	}
	atomic.AddInt64(o.GetCurrentOffset(), t.GetSize())

	t.AddTODOs(
		common.NewFileMD5SumTask,
		common.NewMultipartFileUploadTask,
	)
	return t
}

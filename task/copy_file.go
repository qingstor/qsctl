package task

import (
	"os"
	"sync/atomic"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/pkg/fault"
	"github.com/yunify/qsctl/v2/pkg/types"
	"github.com/yunify/qsctl/v2/task/common"
	"github.com/yunify/qsctl/v2/utils"
)

// newCopyFileTask will create a new copy file task.
func (t *CopyFileTask) new() {
	f, err := os.Open(t.GetPath())
	if err != nil {
		t.TriggerFault(fault.NewUnhandled(err))
		return
	}
	defer f.Close()

	size, err := utils.CalculateFileSize(f)
	if err != nil {
		t.TriggerFault(fault.NewUnhandled(err))
		return
	}
	t.SetTotalSize(size)

	if size >= constants.MaximumAutoMultipartSize {
		t.AddTODOs(NewCopyLargeFileTask)
	} else {
		t.AddTODOs(NewCopySmallFileTask)
	}
	return
}

// newCopySmallFileTask will create a new small file task.
func (t *CopySmallFileTask) new() {
	t.SetOffset(0)
	t.SetSize(t.GetTotalSize())
	t.AddTODOs(
		common.NewFileMD5SumTask,
		common.NewFileUploadTask,
	)
	return
}

// newCopyLargeFileTask will create a new Task.
func (t *CopyLargeFileTask) new() {
	// Init part size.
	partSize, err := utils.CalculatePartSize(t.GetTotalSize())
	if err != nil {
		t.TriggerFault(fault.NewUnhandled(err))
		return
	}
	t.SetPartSize(partSize)

	t.SetScheduler(types.NewScheduler(NewCopyPartialFileTask))

	currentOffset := int64(0)
	t.SetCurrentOffset(&currentOffset)

	t.AddTODOs(
		common.NewMultipartInitTask,
		common.NewWaitTask,
		common.NewMultipartCompleteTask,
	)
	return
}

// NewCopyPartialFileTask will create a new Task.
func (t *CopyPartialFileTask) new() {
	totalSize := t.GetTotalSize()
	partSize := t.GetPartSize()

	// Set part number and update current part number.
	currentOffset := t.GetCurrentOffset()

	// Set size and update offset.
	t.SetOffset(*currentOffset)
	if totalSize < *currentOffset+partSize {
		t.SetSize(totalSize - *currentOffset)
	} else {
		t.SetSize(partSize)
	}
	atomic.AddInt64(t.GetCurrentOffset(), t.GetSize())

	t.AddTODOs(
		common.NewFileMD5SumTask,
		common.NewMultipartFileUploadTask,
	)
	return
}

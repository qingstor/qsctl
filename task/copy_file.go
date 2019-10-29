package task

import (
	"sync/atomic"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/pkg/fault"
	"github.com/yunify/qsctl/v2/task/common"
	"github.com/yunify/qsctl/v2/utils"
)

// newCopyFileTask will create a new copy file task.
func (t *CopyFileTask) new() {
	o, err := t.GetSourceStorage().Stat(t.GetSourcePath())
	if err != nil {
		t.TriggerFault(fault.NewUnhandled(err))
		return
	}

	size, ok := o.GetSize()
	if !ok {
		// TODO: return size not get error.
		t.TriggerFault(fault.NewUnhandled(err))
		return
	}
	t.SetTotalSize(size)
}

func (t *CopyFileTask) run() {
	if t.GetTotalSize() >= constants.MaximumAutoMultipartSize {
		t.GetScheduler().Sync(NewCopyLargeFileTask, t)
	} else {
		t.GetScheduler().Sync(NewCopySmallFileTask, t)
	}
}

// newCopySmallFileTask will create a new small file task.
func (t *CopySmallFileTask) new() {
	t.SetOffset(0)
	t.SetSize(t.GetTotalSize())
}

func (t *CopySmallFileTask) run() {
	t.GetScheduler().Sync(common.NewFileMD5SumTask, t)
	t.GetScheduler().Sync(common.NewFileUploadTask, t)
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

	t.SetScheduleFunc(NewCopyPartialFileTask)

	currentOffset := int64(0)
	t.SetCurrentOffset(&currentOffset)
}

func (t *CopyLargeFileTask) run() {
	t.GetScheduler().Async(common.NewMultipartInitTask, t)
	t.GetScheduler().Sync(common.NewMultipartCompleteTask, t)
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
}

func (t *CopyPartialFileTask) run() {
	t.GetScheduler().Sync(common.NewFileMD5SumTask, t)
	t.GetScheduler().Async(common.NewMultipartFileUploadTask, t)
}

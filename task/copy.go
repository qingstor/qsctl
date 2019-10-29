package task

import (
	"bytes"
	"io"
	"sync"
	"sync/atomic"

	typ "github.com/Xuanwo/storage/types"
	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/pkg/fault"
	"github.com/yunify/qsctl/v2/utils"
)

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
		t.GetScheduler().Sync(t, NewCopyLargeFileTask)
	} else {
		t.GetScheduler().Sync(t, NewCopySmallFileTask)
	}
}

func (t *CopySmallFileTask) new() {
	t.SetOffset(0)
	t.SetSize(t.GetTotalSize())
}

func (t *CopySmallFileTask) run() {
	t.GetScheduler().Sync(t, NewFileMD5SumTask)
	t.GetScheduler().Sync(t, NewFileCopyTask)
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
	t.GetScheduler().Async(t, NewSegmentInitTask)
	t.GetScheduler().Sync(t, NewSegmentCompleteTask)
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
	t.GetScheduler().Sync(t, NewFileMD5SumTask)
	t.GetScheduler().Async(t, NewSegmentFileCopyTask)
}

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

	t.SetScheduleFunc(NewCopyPartialStreamTask)

	currentOffset := int64(0)
	t.SetCurrentOffset(&currentOffset)

	// We don't know how many data in stream, set it to -1 as an indicate.
	// We will set current offset to -1 when got an EOF from stream.
	t.SetTotalSize(-1)
}

func (t *CopyStreamTask) run() {
	t.GetScheduler().Async(t, NewSegmentInitTask)
	t.GetScheduler().Sync(t, NewSegmentCompleteTask)
}

// NewCopyPartialStreamTask will create a new Task.
func (t *CopyPartialStreamTask) new() {
	// Set size and update offset.
	partSize := t.GetPartSize()

	r, err := t.GetSourceStorage().Read(t.GetSourcePath(), typ.WithSize(partSize))
	if err != nil {
		t.TriggerFault(fault.NewUnhandled(err))
		return
	}

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
}

func (t *CopyPartialStreamTask) run() {
	t.GetScheduler().Sync(t, NewStreamMD5SumTask)
	t.GetScheduler().Sync(t, NewSegmentStreamCopyTask)
}

package task

import (
	"bytes"
	"io"
	"sync"

	typ "github.com/Xuanwo/storage/types"
	log "github.com/sirupsen/logrus"
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
}

func (t *CopyLargeFileTask) run() {
	x := newCopyLargeFileShim(t)
	utils.ChooseDestinationStorage(x, t)

	t.GetScheduler().Sync(x, NewSegmentInitTask)
	t.GetScheduler().Sync(x, NewSegmentCompleteTask)
}

// NewCopyPartialFileTask will create a new Task.
func (t *CopyPartialFileTask) new() {
	log.Debugf("Task <%s> for Object <%s> started.", "CopyPartialFile", t.GetDestinationPath())

	totalSize := t.GetTotalSize()
	partSize := t.GetPartSize()
	offset := t.GetOffset()

	if totalSize < offset+partSize {
		t.SetSize(totalSize - offset)
		t.SetDone(true)
	} else {
		t.SetSize(partSize)
	}

	log.Debugf("Task <%s> for Object <%s> started.", "CopyPartialFile", t.GetDestinationPath())
}

func (t *CopyPartialFileTask) run() {
	t.GetScheduler().Sync(t, NewFileMD5SumTask)
	t.GetScheduler().Sync(t, NewSegmentFileCopyTask)
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
	log.Debugf("Task <%s> for Object <%s> started.", "CopyPartialStream", t.GetDestinationPath())

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

	t.SetSize(n)
	t.SetContent(b)
	if n < partSize {
		t.SetDone(true)
	}

	log.Debugf("Task <%s> for Object <%s> finished.", "CopyPartialStream", t.GetDestinationPath())
}

func (t *CopyPartialStreamTask) run() {
	t.GetScheduler().Sync(t, NewStreamMD5SumTask)
	t.GetScheduler().Sync(t, NewSegmentStreamCopyTask)
}

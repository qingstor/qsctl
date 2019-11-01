package task

import (
	"bytes"
	"io"
	"sync"

	typ "github.com/Xuanwo/storage/types"
	log "github.com/sirupsen/logrus"
	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/pkg/types"
	"github.com/yunify/qsctl/v2/utils"
)

func (t *CopyFileTask) new() {
	o, err := t.GetSourceStorage().Stat(t.GetSourcePath())
	if err != nil {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}

	size, ok := o.GetSize()
	if !ok {
		// TODO: return size not get error.
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}
	t.SetTotalSize(size)
}

func (t *CopyFileTask) run() {
	if t.GetTotalSize() >= constants.MaximumAutoMultipartSize {
		t.GetScheduler().Sync(NewCopyLargeFileTask(t))
	} else {
		t.GetScheduler().Sync(NewCopySmallFileTask(t))
	}
}

func (t *CopySmallFileTask) new() {
	t.SetOffset(0)
	t.SetSize(t.GetTotalSize())
}

func (t *CopySmallFileTask) run() {
	t.GetScheduler().Sync(NewMD5SumFileTask(t))
	t.GetScheduler().Sync(NewCopySingleFileTask(t))
}

// newCopyLargeFileTask will create a new Task.
func (t *CopyLargeFileTask) new() {
}

func (t *CopyLargeFileTask) run() {
	initTask := NewSegmentInit(t)
	utils.ChooseDestinationStorage(initTask, t)

	// Set segment part size.
	partSize, err := utils.CalculatePartSize(t.GetTotalSize())
	if err != nil {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}
	initTask.SetPartSize(partSize)
	initTask.SetSegmentScheduleFunc(NewCopyPartialFileSegmentRequirement)

	t.GetScheduler().Sync(initTask)
	t.GetScheduler().Sync(NewSegmentCompleteTask(initTask))
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
		t.SetDone(false)
	}

	log.Debugf("Task <%s> for Object <%s> started.", "CopyPartialFile", t.GetDestinationPath())
}

func (t *CopyPartialFileTask) run() {
	t.GetScheduler().Sync(NewMD5SumFileTask(t))
	t.GetScheduler().Sync(NewSegmentFileCopyTask(t))
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
	t.GetScheduler().Async(NewSegmentInitTask(t))
	t.GetScheduler().Sync(NewSegmentCompleteTask(t))
}

// NewCopyPartialStreamTask will create a new Task.
func (t *CopyPartialStreamTask) new() {
	log.Debugf("Task <%s> for Object <%s> started.", "CopyPartialStream", t.GetDestinationPath())

	// Set size and update offset.
	partSize := t.GetPartSize()

	r, err := t.GetSourceStorage().Read(t.GetSourcePath(), typ.WithSize(partSize))
	if err != nil {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}

	b := t.GetBytesPool().Get().(*bytes.Buffer)
	n, err := io.Copy(b, r)
	if err != nil {
		t.TriggerFault(types.NewErrUnhandled(err))
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
	t.GetScheduler().Sync(NewMD5SumStreamTask(t))
	t.GetScheduler().Sync(NewSegmentStreamCopyTask(t))
}

func (t *CopySingleFileTask) new() {}
func (t *CopySingleFileTask) run() {
	log.Debugf("Task <%s> for file from <%s> to <%s> started.", "FileCopy", t.GetSourcePath(), t.GetDestinationPath())

	r, err := t.GetSourceStorage().Read(t.GetSourcePath())
	if err != nil {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}
	defer r.Close()

	// TODO: add checksum support
	err = t.GetDestinationStorage().Write(t.GetDestinationPath(), r, typ.WithSize(t.GetSize()))
	if err != nil {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}

	log.Debugf("Task <%s> for file from <%s> to <%s> started.", "FileUpload", t.GetSourcePath(), t.GetDestinationPath())
}

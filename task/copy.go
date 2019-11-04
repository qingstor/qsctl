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
	md5Task := NewMD5SumFile(t)
	t.GetScheduler().Sync(md5Task)

	fileCopyTask := NewCopySingleFile(t)
	fileCopyTask.SetMD5Sum(md5Task.GetMD5Sum())
	t.GetScheduler().Sync(fileCopyTask)
}

// newCopyLargeFileTask will create a new Task.
func (t *CopyLargeFileTask) new() {
}

func (t *CopyLargeFileTask) run() {
	// Set segment part size.
	partSize, err := utils.CalculatePartSize(t.GetTotalSize())
	if err != nil {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}
	t.SetPartSize(partSize)

	initTask := NewSegmentInit(t)
	utils.ChooseDestinationStorage(initTask, t)

	t.GetScheduler().Sync(initTask)
	t.SetSegmentID(initTask.GetSegmentID())

	offset := int64(0)
	for {
		t.SetOffset(offset)

		x := NewCopyPartialFile(t)
		t.GetScheduler().Async(x)
		// While GetDone is true, this must be the last part.
		if x.GetDone() {
			break
		}

		offset += x.GetSize()
	}

	// Make sure all segment upload finished.
	t.GetScheduler().Wait()
	t.GetScheduler().Sync(NewSegmentCompleteTask(initTask))
}

// NewCopyPartialFileTask will create a new Task.
func (t *CopyPartialFileTask) new() {
	log.Debugf("Task <%s> for Object <%s> started.", "CopyPartialFile", t.GetDestinationPath())

	totalSize := t.GetTotalSize()
	partSize := t.GetPartSize()
	offset := t.GetOffset()

	if totalSize <= offset+partSize {
		t.SetSize(totalSize - offset)
		t.SetDone(true)
	} else {
		t.SetSize(partSize)
		t.SetDone(false)
	}

	log.Debugf("Task <%s> for Object <%s> finished.", "CopyPartialFile", t.GetDestinationPath())
}

func (t *CopyPartialFileTask) run() {
	md5Task := NewMD5SumFile(t)
	t.GetScheduler().Sync(md5Task)

	fileCopyTask := NewSegmentFileCopy(t)
	fileCopyTask.SetMD5Sum(md5Task.GetMD5Sum())
	t.GetScheduler().Sync(fileCopyTask)
}

// NewCopyStreamTask will create a copy stream task.
func (t *CopyStreamTask) new() {
	bytesPool := &sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 0, t.GetPartSize()))
		},
	}
	t.SetBytesPool(bytesPool)
}

func (t *CopyStreamTask) run() {
	initTask := NewSegmentInit(t)
	utils.ChooseDestinationStorage(initTask, t)

	// TODO: we will use expect size to calculate part size later.
	partSize := int64(constants.DefaultPartSize)
	t.SetPartSize(partSize)

	t.GetScheduler().Sync(initTask)
	t.SetSegmentID(initTask.GetSegmentID())

	for {
		x := NewCopyPartialStream(t)
		t.GetScheduler().Async(x)

		if x.GetDone() {
			break
		}
	}

	t.GetScheduler().Wait()
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

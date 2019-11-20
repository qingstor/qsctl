package task

import (
	"bytes"
	"io"
	"sync"

	typ "github.com/Xuanwo/storage/types"
	"github.com/Xuanwo/storage/types/pairs"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/pkg/types"
	"github.com/yunify/qsctl/v2/utils"
)

func (t *CopyDirTask) new() {}
func (t *CopyDirTask) run() {
	x := NewListDir(t)
	utils.ChooseSourceStorage(x, t)
	x.SetFileFunc(func(o *typ.Object) {
		sf := NewCopyFile(t)
		sf.SetSourcePath(o.Name)
		sf.SetDestinationPath(o.Name)
		t.GetScheduler().Async(sf)
	})
	x.SetDirFunc(func(o *typ.Object) {
		sf := NewCopyDir(t)
		sf.SetSourcePath(o.Name)
		sf.SetDestinationPath(o.Name)
		t.GetScheduler().Sync(sf)
	})
	t.GetScheduler().Sync(x)
}

func (t *CopyFileTask) new() {}
func (t *CopyFileTask) run() {
	o, err := t.GetSourceStorage().Stat(t.GetSourcePath())
	if err != nil {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}

	if o.Size >= constants.MaximumAutoMultipartSize {
		x := NewCopyLargeFile(t)
		x.SetTotalSize(o.Size)
		t.GetScheduler().Sync(x)
	} else {
		x := NewCopySmallFile(t)
		x.SetSize(o.Size)
		t.GetScheduler().Sync(x)
	}
}

func (t *CopySmallFileTask) new() {}
func (t *CopySmallFileTask) run() {
	md5Task := NewMD5SumFile(t)
	utils.ChooseSourceStorage(md5Task, t)
	md5Task.SetOffset(0)
	t.GetScheduler().Sync(md5Task)

	fileCopyTask := NewCopySingleFile(t)
	fileCopyTask.SetMD5Sum(md5Task.GetMD5Sum())
	t.GetScheduler().Sync(fileCopyTask)
}

func (t *CopyLargeFileTask) new() {}
func (t *CopyLargeFileTask) run() {
	// Set segment part size.
	partSize, err := utils.CalculatePartSize(t.GetTotalSize())
	if err != nil {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}
	t.SetPartSize(partSize)

	initTask := NewSegmentInit(t)
	err = utils.ChooseDestinationStorageAsSegmenter(initTask, t)
	if err != nil {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}

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

func (t *CopyPartialFileTask) new() {
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
}
func (t *CopyPartialFileTask) run() {
	md5Task := NewMD5SumFile(t)
	utils.ChooseSourceStorage(md5Task, t)
	t.GetScheduler().Sync(md5Task)

	fileCopyTask := NewSegmentFileCopy(t)
	fileCopyTask.SetMD5Sum(md5Task.GetMD5Sum())
	err := utils.ChooseDestinationSegmenter(fileCopyTask, t)
	if err != nil {
		t.TriggerFault(err)
		return
	}
	t.GetScheduler().Sync(fileCopyTask)
}

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
	err := utils.ChooseDestinationStorageAsSegmenter(initTask, t)
	if err != nil {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}

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

func (t *CopyPartialStreamTask) new() {
	// Set size and update offset.
	partSize := t.GetPartSize()

	r, err := t.GetSourceStorage().Read(t.GetSourcePath(), pairs.WithSize(partSize))
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
	} else {
		t.SetDone(false)
	}
}
func (t *CopyPartialStreamTask) run() {
	t.GetScheduler().Sync(NewMD5SumStreamTask(t))
	t.GetScheduler().Sync(NewSegmentStreamCopyTask(t))
}

func (t *CopySingleFileTask) new() {}
func (t *CopySingleFileTask) run() {
	r, err := t.GetSourceStorage().Read(t.GetSourcePath())
	if err != nil {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}
	defer r.Close()

	// TODO: add checksum support
	err = t.GetDestinationStorage().Write(t.GetDestinationPath(), r, pairs.WithSize(t.GetSize()))
	if err != nil {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}
}

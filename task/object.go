package task

import (
	"bufio"
	"io"
	"os"
	"sync"

	"github.com/Xuanwo/navvy"
	log "github.com/sirupsen/logrus"

	"github.com/yunify/qsctl/v2/task/types"
	"github.com/yunify/qsctl/v2/task/utils"
)

type PutObjectTaskRequirement interface {
	navvy.Task
	types.Todoist

	types.ObjectKeyGetter
	types.FilePathGetter
	types.MD5SumGetter

	types.StorageGetter
	types.PoolGetter
}

type PutObjectTask struct {
	PutObjectTaskRequirement
}

// NewPutObjectTask will create a new Task.
func NewPutObjectTask(task types.Todoist) navvy.Task {
	o, ok := task.(PutObjectTaskRequirement)
	if !ok {
		panic("task is not fill PutObjectTaskRequirement")
	}

	return &PutObjectTask{o}
}

// Run implement navvy.Task.
func (t *PutObjectTask) Run() {
	log.Debugf("Task <%s> for Object <%s> started.", "PutObjectTask", t.GetObjectKey())

	f, err := os.Open(t.GetFilePath())
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = t.GetStorage().PutObject(t.GetObjectKey(), t.GetMD5Sum(), f)
	if err != nil {
		panic(err)
	}

	log.Debugf("Task <%s> for Object <%s> finished.", "PutObjectTask", t.GetObjectKey())
	utils.SubmitNextTask(t.PutObjectTaskRequirement)
}

// MultipartObjectInitTaskRequirement is the requirement for execute MultipartObjectInitTask.
type MultipartObjectInitTaskRequirement interface {
	navvy.Task

	types.Todoist

	types.ObjectKeyGetter
	types.FilePathGetter

	types.UploadIDSetter
	types.TotalPartsSetter
	types.WaitGroupSetter
	types.PoolGetter
	types.StorageGetter
}

// MultipartObjectInitTask will execute MultipartObjectInit Task.
type MultipartObjectInitTask struct {
	MultipartObjectInitTaskRequirement
}

// NewMultipartObjectInitTask will create a new Task.
func NewMultipartObjectInitTask(task types.Todoist) navvy.Task {
	o, ok := task.(MultipartObjectInitTaskRequirement)
	if !ok {
		panic("task is not fill MultipartObjectInitTaskRequirement")
	}

	return &MultipartObjectInitTask{o}
}

// Run implement navvy.Task.
func (t *MultipartObjectInitTask) Run() {
	log.Debugf("Task <%s> for Object <%s> started.", "MultipartObjectInitTask", t.GetObjectKey())

	f, err := os.Open(t.GetFilePath())
	if err != nil {
		panic(err)
	}
	defer f.Close()

	uploadID, err := t.GetStorage().InitiateMultipartUpload(t.GetObjectKey())
	if err != nil {
		panic(err)
	}
	t.SetUploadID(uploadID)

	size, err := utils.CalculateSeekableFileSize(f)
	if err != nil {
		panic(err)
	}

	partSize, err := utils.CalculatePartSize(size)
	if err != nil {
		panic(err)
	}

	wg := &sync.WaitGroup{}
	t.SetWaitGroup(wg)

	partNumber := 0
	cur := int64(0)
	for {
		curPartSize := partSize
		if cur+partSize > size {
			curPartSize = size - cur
		}

		task := NewCopyPartialFileTask(
			t.GetObjectKey(),
			t.GetFilePath(),
			uploadID,
			partNumber,
			cur,
			curPartSize,
		)
		wg.Add(1)
		task.SetWaitGroup(wg)

		log.Debugf("Submit task <%s> for Object <%s> with UploadID <%s> in PartNumber <%d>", "CopyPartialFileTask",
			task.GetObjectKey(), task.GetUploadID(), task.GetPartNumber())
		go t.GetPool().Submit(task)

		partNumber++
		cur += curPartSize
		if cur >= size {
			break
		}
	}

	t.SetTotalParts(partNumber)

	log.Debugf("Task <%s> for Object <%s> finished.", "MultipartObjectInitTask", t.GetObjectKey())
	utils.SubmitNextTask(t.MultipartObjectInitTaskRequirement)
}

// MultipartObjectUploadTaskRequirement is the requirement for execute MultipartObjectUploadTask.
type MultipartObjectUploadTaskRequirement interface {
	navvy.Task
	types.Todoist

	types.MD5SumGetter
	types.FilePathGetter
	types.ObjectKeyGetter
	types.OffsetGetter
	types.UploadIDGetter
	types.PartNumberGetter
	types.SizeGetter
	types.WaitGroupGetter
	types.PoolGetter
	types.StorageGetter
}

// MultipartObjectUploadTask will execute MultipartObjectUpload Task.
type MultipartObjectUploadTask struct {
	MultipartObjectUploadTaskRequirement
}

// NewMultipartObjectUploadTask will create a new Task.
func NewMultipartObjectUploadTask(task types.Todoist) navvy.Task {
	o, ok := task.(MultipartObjectUploadTaskRequirement)
	if !ok {
		panic("task is not fill MultipartObjectUploadTaskRequirement")
	}

	return &MultipartObjectUploadTask{o}
}

// Run implement navvy.Task.
func (t *MultipartObjectUploadTask) Run() {
	log.Debugf("Task <%s> for File <%s> at Offset <%d> started.", "MultipartObjectUploadTask", t.GetFilePath(), t.GetOffset())

	f, err := os.Open(t.GetFilePath())
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := bufio.NewReader(io.NewSectionReader(f, t.GetOffset(), t.GetSize()))

	err = t.GetStorage().UploadMultipart(t.GetObjectKey(), t.GetUploadID(), t.GetSize(), t.GetPartNumber(), t.GetMD5Sum(), r)
	if err != nil {
		panic(err)
	}

	t.GetWaitGroup().Done()

	log.Debugf("Task <%s> for File <%s> at Offset <%d> finished.", "MultipartObjectUploadTask", t.GetFilePath(), t.GetOffset())
	utils.SubmitNextTask(t.MultipartObjectUploadTaskRequirement)
}

// MultipartObjectCompleteTaskRequirement will execute MultipartObjectCompleteT Task.
type MultipartObjectCompleteTaskRequirement interface {
	navvy.Task
	types.Todoist

	types.ObjectKeyGetter
	types.UploadIDGetter
	types.TotalPartsGetter
	types.PoolGetter
	types.StorageGetter
}

// MultipartObjectCompleteTask will execute MultipartObjectComplete Task.
type MultipartObjectCompleteTask struct {
	MultipartObjectCompleteTaskRequirement
}

// NewMultipartObjectCompleteTask will create a new Task.
func NewMultipartObjectCompleteTask(task types.Todoist) navvy.Task {
	o, ok := task.(MultipartObjectCompleteTaskRequirement)
	if !ok {
		panic("task is not fill NewMultipartObjectCompleteTask")
	}

	return &MultipartObjectCompleteTask{o}
}

// Run implement navvy.Task.
func (t *MultipartObjectCompleteTask) Run() {
	log.Debugf("Task <%s> for Object <%s> UploadID <%s> started.", "MultipartObjectCompleteTask", t.GetObjectKey(), t.GetUploadID())

	err := t.GetStorage().CompleteMultipartUpload(t.GetObjectKey(), t.GetUploadID(), t.GetTotalParts())
	if err != nil {
		panic(err)
	}

	log.Debugf("Task <%s> for Object <%s> UploadID <%s> finished.", "MultipartObjectCompleteTask", t.GetObjectKey(), t.GetUploadID())
	utils.SubmitNextTask(t.MultipartObjectCompleteTaskRequirement)
}
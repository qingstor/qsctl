package task

import (
	"bufio"
	"io"
	"os"
	"sync"

	"github.com/Xuanwo/navvy"
	log "github.com/sirupsen/logrus"
	"github.com/yunify/qsctl/v2/contexts"
)

// MultipartObjectInitTaskRequirement is the requirement for execute MultipartObjectInitTask.
type MultipartObjectInitTaskRequirement interface {
	Todoist

	ObjectKeyGetter
	FilePathGetter

	UploadIDSetter
	TotalPartsSetter
	WaitGroupSetter
}

// MultipartObjectInitTask will execute MultipartObjectInit Task.
type MultipartObjectInitTask struct {
	MultipartObjectInitTaskRequirement
}

// NewMultipartObjectInitTask will create a new Task.
func NewMultipartObjectInitTask(task Todoist) navvy.Task {
	o, ok := task.(MultipartObjectInitTaskRequirement)
	if !ok {
		panic("task is not fill MultipartObjectInitTaskRequirement")
	}

	return &MultipartObjectInitTask{o}
}

// Run implement navvy.Task.
func (t *MultipartObjectInitTask) Run() {
	log.Debugf("Task <%s> for Object <%s> started.", "MultipartObjectInitTask", t.GetObjectKey())

	// TODO: check file size.
	f, err := os.Open(t.GetFilePath())
	if err != nil {
		panic(err)
	}
	defer f.Close()

	uploadID, err := contexts.Storage.InitiateMultipartUpload(t.GetObjectKey())
	if err != nil {
		panic(err)
	}
	t.SetUploadID(uploadID)

	size, err := CalculateSeekableFileSize(f)
	if err != nil {
		panic(err)
	}

	partSize, err := CalculatePartSize(size)
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
		go contexts.Pool.Submit(task)

		partNumber++
		cur += curPartSize
		if cur >= size {
			break
		}
	}

	t.SetTotalParts(partNumber)

	log.Debugf("Task <%s> for Object <%s> finished.", "MultipartObjectInitTask", t.GetObjectKey())
	go SubmitNextTask(t.MultipartObjectInitTaskRequirement)
}

// MultipartObjectUploadTaskRequirement is the requirement for execute MultipartObjectUploadTask.
type MultipartObjectUploadTaskRequirement interface {
	navvy.Task
	Todoist

	MD5SumGetter
	FilePathGetter
	ObjectKeyGetter
	OffsetGetter
	UploadIDGetter
	PartNumberGetter
	ContentLengthGetter
	WaitGroupGetter
}

// MultipartObjectUploadTask will execute MultipartObjectUpload Task.
type MultipartObjectUploadTask struct {
	MultipartObjectUploadTaskRequirement
}

// NewMultipartObjectUploadTask will create a new Task.
func NewMultipartObjectUploadTask(task Todoist) navvy.Task {
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

	r := bufio.NewReader(io.NewSectionReader(f, t.GetOffset(), t.GetContentLength()))

	err = contexts.Storage.UploadMultipart(t.GetObjectKey(), t.GetUploadID(), t.GetContentLength(), t.GetPartNumber(), t.GetMD5Sum(), r)
	if err != nil {
		panic(err)
	}

	t.GetWaitGroup().Done()

	log.Debugf("Task <%s> for File <%s> at Offset <%d> finished.", "MultipartObjectUploadTask", t.GetFilePath(), t.GetOffset())
	go SubmitNextTask(t.MultipartObjectUploadTaskRequirement)
}

// MultipartObjectCompleteTaskRequirement will execute MultipartObjectCompleteT Task.
type MultipartObjectCompleteTaskRequirement interface {
	navvy.Task
	Todoist

	ObjectKeyGetter
	UploadIDGetter
	TotalPartsGetter
}

// MultipartObjectCompleteTask will execute MultipartObjectComplete Task.
type MultipartObjectCompleteTask struct {
	MultipartObjectCompleteTaskRequirement
}

// NewMultipartObjectCompleteTask will create a new Task.
func NewMultipartObjectCompleteTask(task Todoist) navvy.Task {
	o, ok := task.(MultipartObjectCompleteTaskRequirement)
	if !ok {
		panic("task is not fill NewMultipartObjectCompleteTask")
	}

	return &MultipartObjectCompleteTask{o}
}

// Run implement navvy.Task.
func (t *MultipartObjectCompleteTask) Run() {
	log.Debugf("Task <%s> for Object <%s> UploadID <%s> started.", "MultipartObjectCompleteTask", t.GetObjectKey(), t.GetUploadID())

	err := contexts.Storage.CompleteMultipartUpload(t.GetObjectKey(), t.GetUploadID(), t.GetTotalParts())
	if err != nil {
		panic(err)
	}

	log.Debugf("Task <%s> for Object <%s> UploadID <%s> finished.", "MultipartObjectCompleteTask", t.GetObjectKey(), t.GetUploadID())
	go SubmitNextTask(t.MultipartObjectCompleteTaskRequirement)
}

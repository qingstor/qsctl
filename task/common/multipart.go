package common

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

// MultipartInitTaskRequirement is the requirement for execute MultipartInitTask.
type MultipartInitTaskRequirement interface {
	navvy.Task

	types.Todoist

	types.ObjectKeyGetter
	types.FilePathGetter

	types.UploadIDSetter
	types.WaitGroupSetter
	types.PoolGetter
	types.StorageGetter
	types.TaskConstructorGetter
	types.CurrentPartNumberGetter
	types.CurrentOffsetGetter
	types.PartSizeGetter
	types.SizeGetter
}

// MultipartInitTask will execute MultipartObjectInit Task.
type MultipartInitTask struct {
	MultipartInitTaskRequirement
}

// NewMultipartInitTask will create a new Task.
func NewMultipartInitTask(task types.Todoist) navvy.Task {
	o, ok := task.(MultipartInitTaskRequirement)
	if !ok {
		panic("task is not fill MultipartInitTaskRequirement")
	}

	return &MultipartInitTask{o}
}

// Run implement navvy.Task.
func (t *MultipartInitTask) Run() {
	log.Debugf("Task <%s> for Object <%s> started.", "MultipartInitTask", t.GetObjectKey())

	uploadID, err := t.GetStorage().InitiateMultipartUpload(t.GetObjectKey())
	if err != nil {
		panic(err)
	}
	t.SetUploadID(uploadID)

	wg := &sync.WaitGroup{}
	t.SetWaitGroup(wg)

	for {
		if *t.GetCurrentOffset() == t.GetSize() {
			break
		}

		task := t.GetTaskConstructor()(t)
		wg.Add(1)

		go t.GetPool().Submit(task)
	}

	log.Debugf("Task <%s> for Object <%s> finished.", "MultipartInitTask", t.GetObjectKey())
	utils.SubmitNextTask(t.MultipartInitTaskRequirement)
}

// MultipartFileUploadTaskRequirement is the requirement for execute MultipartFileUploadTask.
type MultipartFileUploadTaskRequirement interface {
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

// MultipartFileUploadTask will execute MultipartObjectUpload Task.
type MultipartFileUploadTask struct {
	MultipartFileUploadTaskRequirement
}

// NewMultipartFileUploadTask will create a new Task.
func NewMultipartFileUploadTask(task types.Todoist) navvy.Task {
	o, ok := task.(MultipartFileUploadTaskRequirement)
	if !ok {
		panic("task is not fill MultipartFileUploadTaskRequirement")
	}

	return &MultipartFileUploadTask{o}
}

// Run implement navvy.Task.
func (t *MultipartFileUploadTask) Run() {
	log.Debugf("Task <%s> for File <%s> at Offset <%d> started.", "MultipartFileUploadTask", t.GetFilePath(), t.GetOffset())

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

	log.Debugf("Task <%s> for File <%s> at Offset <%d> finished.", "MultipartFileUploadTask", t.GetFilePath(), t.GetOffset())
	utils.SubmitNextTask(t.MultipartFileUploadTaskRequirement)
}

// MultipartStreamUploadTaskRequirement is the requirement for execute MultipartStreamUploadTask.
type MultipartStreamUploadTaskRequirement interface {
	navvy.Task
	types.Todoist

	types.MD5SumGetter
	types.FilePathGetter
	types.ObjectKeyGetter
	types.UploadIDGetter
	types.PartNumberGetter
	types.SizeGetter
	types.WaitGroupGetter
	types.PoolGetter
	types.StorageGetter
	types.ContentGetter
}

// MultipartStreamUploadTask will execute MultipartObjectUpload Task.
type MultipartStreamUploadTask struct {
	MultipartStreamUploadTaskRequirement
}

// NewMultipartStreamUploadTask will create a new Task.
func NewMultipartStreamUploadTask(task types.Todoist) navvy.Task {
	o, ok := task.(MultipartStreamUploadTaskRequirement)
	if !ok {
		panic("task is not fill MultipartStreamUploadTaskRequirement")
	}

	return &MultipartStreamUploadTask{o}
}

// Run implement navvy.Task.
func (t *MultipartStreamUploadTask) Run() {
	log.Debugf("Task <%s> for Stream <%s> at PartNumber <%d> started.", "MultipartStreamUploadTask", t.GetFilePath(), t.GetPartNumber())

	err := t.GetStorage().UploadMultipart(
		t.GetObjectKey(), t.GetUploadID(), t.GetSize(),
		t.GetPartNumber(), t.GetMD5Sum(), t.GetContent())
	if err != nil {
		panic(err)
	}

	t.GetWaitGroup().Done()

	log.Debugf("Task <%s> for Stream <%s> at PartNumber <%d> finished.", "MultipartStreamUploadTask", t.GetFilePath(), t.GetPartNumber())
	utils.SubmitNextTask(t.MultipartStreamUploadTaskRequirement)
}

// MultipartCompleteTaskRequirement will execute MultipartObjectCompleteT Task.
type MultipartCompleteTaskRequirement interface {
	navvy.Task
	types.Todoist

	types.ObjectKeyGetter
	types.UploadIDGetter
	types.CurrentPartNumberGetter
	types.PoolGetter
	types.StorageGetter
}

// MultipartCompleteTask will execute MultipartObjectComplete Task.
type MultipartCompleteTask struct {
	MultipartCompleteTaskRequirement
}

// NewMultipartCompleteTask will create a new Task.
func NewMultipartCompleteTask(task types.Todoist) navvy.Task {
	o, ok := task.(MultipartCompleteTaskRequirement)
	if !ok {
		panic("task is not fill NewMultipartCompleteTask")
	}

	return &MultipartCompleteTask{o}
}

// Run implement navvy.Task.
func (t *MultipartCompleteTask) Run() {
	log.Debugf("Task <%s> for Object <%s> UploadID <%s> started.", "MultipartCompleteTask", t.GetObjectKey(), t.GetUploadID())

	err := t.GetStorage().CompleteMultipartUpload(t.GetObjectKey(), t.GetUploadID(), int(*t.GetCurrentPartNumber()-1))
	if err != nil {
		panic(err)
	}

	log.Debugf("Task <%s> for Object <%s> UploadID <%s> finished.", "MultipartCompleteTask", t.GetObjectKey(), t.GetUploadID())
	utils.SubmitNextTask(t.MultipartCompleteTaskRequirement)
}

package common

import (
	"bufio"
	"io"
	"os"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/yunify/qsctl/v2/task/utils"
)

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

// Run implement navvy.Task.
func (t *MultipartFileUploadTask) Run() {
	log.Debugf("Task <%s> for File <%s> at Offset <%d> started.", "MultipartFileUploadTask", t.GetPath(), t.GetOffset())

	f, err := os.Open(t.GetPath())
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

	log.Debugf("Task <%s> for File <%s> at Offset <%d> finished.", "MultipartFileUploadTask", t.GetPath(), t.GetOffset())
	utils.SubmitNextTask(t.MultipartFileUploadTaskRequirement)
}

// Run implement navvy.Task.
func (t *MultipartStreamUploadTask) Run() {
	log.Debugf("Task <%s> for Stream <%s> at PartNumber <%d> started.", "MultipartStreamUploadTask", t.GetPath(), t.GetPartNumber())

	err := t.GetStorage().UploadMultipart(
		t.GetObjectKey(), t.GetUploadID(), t.GetSize(),
		t.GetPartNumber(), t.GetMD5Sum(), t.GetContent())
	if err != nil {
		panic(err)
	}

	t.GetWaitGroup().Done()

	log.Debugf("Task <%s> for Stream <%s> at PartNumber <%d> finished.", "MultipartStreamUploadTask", t.GetPath(), t.GetPartNumber())
	utils.SubmitNextTask(t.MultipartStreamUploadTaskRequirement)
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

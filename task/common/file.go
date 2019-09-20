package common

import (
	"os"

	"github.com/Xuanwo/navvy"
	log "github.com/sirupsen/logrus"

	"github.com/yunify/qsctl/v2/task/types"
	"github.com/yunify/qsctl/v2/task/utils"
)

type FileUploadTaskRequirement interface {
	navvy.Task
	types.Todoist

	types.ObjectKeyGetter
	types.PathGetter
	types.MD5SumGetter

	types.StorageGetter
	types.PoolGetter
}

type FileUploadTask struct {
	FileUploadTaskRequirement
}

// NewFileUploadTask will create a new Task.
func NewFileUploadTask(task types.Todoist) navvy.Task {
	o, ok := task.(FileUploadTaskRequirement)
	if !ok {
		panic("task is not fill FileUploadTaskRequirement")
	}

	return &FileUploadTask{o}
}

// Run implement navvy.Task.
func (t *FileUploadTask) Run() {
	log.Debugf("Task <%s> for Object <%s> started.", "FileUploadTask", t.GetObjectKey())

	f, err := os.Open(t.GetPath())
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = t.GetStorage().PutObject(t.GetObjectKey(), t.GetMD5Sum(), f)
	if err != nil {
		panic(err)
	}

	log.Debugf("Task <%s> for Object <%s> finished.", "FileUploadTask", t.GetObjectKey())
	utils.SubmitNextTask(t.FileUploadTaskRequirement)
}

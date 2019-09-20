package common

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/yunify/qsctl/v2/task/utils"
)

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

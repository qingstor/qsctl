package common

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func (t *FileUploadTask) run() {
	log.Debugf("Task <%s> for Object <%s> started.", "FileUploadTask", t.GetKey())

	f, err := os.Open(t.GetPath())
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = t.GetStorage().PutObject(t.GetKey(), t.GetMD5Sum(), f)
	if err != nil {
		panic(err)
	}

	log.Debugf("Task <%s> for Object <%s> finished.", "FileUploadTask", t.GetKey())
}

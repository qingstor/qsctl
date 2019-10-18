package common

import (
	"os"

	"github.com/Xuanwo/storage/types"
	log "github.com/sirupsen/logrus"

	"github.com/yunify/qsctl/v2/pkg/fault"
)

func (t *FileUploadTask) run() {
	log.Debugf("Task <%s> for Object <%s> started.", "FileUploadTask", t.GetKey())

	f, err := os.Open(t.GetPath())
	if err != nil {
		t.TriggerFault(fault.NewUnhandled(err))
		return
	}
	defer f.Close()

	// TODO: add checksum support
	err = t.GetDestinationStorage().Write(t.GetKey(), f, types.WithSize(t.GetSize()))
	if err != nil {
		panic(err)
	}

	log.Debugf("Task <%s> for Object <%s> finished.", "FileUploadTask", t.GetKey())
}

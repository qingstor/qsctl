package task

import (
	typ "github.com/Xuanwo/storage/types"
	log "github.com/sirupsen/logrus"

	"github.com/yunify/qsctl/v2/pkg/fault"
)

func (t *FileUploadTask) run() {
	log.Debugf("Task <%s> for Object <%s> started.", "FileUploadTask", t.GetDestinationPath())

	r, err := t.GetSourceStorage().Read(t.GetSourcePath())
	if err != nil {
		t.TriggerFault(fault.NewUnhandled(err))
		return
	}
	defer r.Close()

	// TODO: add checksum support
	err = t.GetDestinationStorage().Write(t.GetDestinationPath(), r, typ.WithSize(t.GetSize()))
	if err != nil {
		panic(err)
	}

	log.Debugf("Task <%s> for Object <%s> finished.", "FileUploadTask", t.GetDestinationPath())
}

package task

import (
	typ "github.com/Xuanwo/storage/types"
	log "github.com/sirupsen/logrus"

	"github.com/yunify/qsctl/v2/pkg/fault"
)

func (t *FileCopyTask) new() {}
func (t *FileCopyTask) run() {
	log.Debugf("Task <%s> for file from <%s> to <%s> started.", "FileCopy", t.GetSourcePath(), t.GetDestinationPath())

	r, err := t.GetSourceStorage().Read(t.GetSourcePath())
	if err != nil {
		t.TriggerFault(fault.NewUnhandled(err))
		return
	}
	defer r.Close()

	// TODO: add checksum support
	err = t.GetDestinationStorage().Write(t.GetDestinationPath(), r, typ.WithSize(t.GetSize()))
	if err != nil {
		t.TriggerFault(fault.NewUnhandled(err))
		return
	}

	log.Debugf("Task <%s> for file from <%s> to <%s> started.", "FileUpload", t.GetSourcePath(), t.GetDestinationPath())
}

package task

import (
	"github.com/Xuanwo/navvy"
	typ "github.com/Xuanwo/storage/types"
	log "github.com/sirupsen/logrus"

	"github.com/yunify/qsctl/v2/pkg/fault"
	"github.com/yunify/qsctl/v2/pkg/types"
)

func (t *DeleteDirTask) new() {
	oc := make(chan *typ.Object)
	t.SetObjectChannel(oc)

	// set recursive for list async task to list recursively
	t.SetRecursive(true)

	t.SetScheduleFunc(NewDeleteFileTask)
}

func (t *DeleteDirTask) run() {
	log.Debugf("Task <%s> for path <%s> started",
		"DeleteDir", t.GetPath())

	t.GetScheduler().Sync(t, NewFileIteratorTask)

	log.Debugf("Task <%s> for path <%s> finished",
		"DeleteDir", t.GetPath())
}

func (t *DeleteFileTask) new() {}

func (t *DeleteFileTask) run() {
	log.Debugf("Task <%s> for path <%s> started",
		"DeleteFile", t.GetPath())

	if err := t.GetStorage().Delete(t.GetPath()); err != nil {
		t.TriggerFault(fault.NewUnhandled(err))
		return
	}

	log.Debugf("Task <%s> for path <%s> finished",
		"DeleteFile", t.GetPath())
}

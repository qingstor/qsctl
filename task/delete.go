package task

import (
	log "github.com/sirupsen/logrus"
	"github.com/yunify/qsctl/v2/pkg/types"
)

func (t *DeleteDirTask) new() {}

func (t *DeleteDirTask) run() {
	log.Debugf("Task <%s> for path <%s> started",
		"DeleteDir", t.GetPath())

	x := NewIterateFile(t)
	x.SetPathScheduleFunc(NewDeleteFilePathRequirement)
	x.SetRecursive(true)
	t.GetScheduler().Sync(x)

	log.Debugf("Task <%s> for path <%s> finished",
		"DeleteDir", t.GetPath())
}

func (t *DeleteFileTask) new() {}

func (t *DeleteFileTask) run() {
	log.Debugf("Task <%s> for path <%s> started",
		"DeleteFile", t.GetPath())

	if err := t.GetStorage().Delete(t.GetPath()); err != nil {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}

	log.Debugf("Task <%s> for path <%s> finished",
		"DeleteFile", t.GetPath())
}

func (t *DeleteStorageTask) new() {
}

func (t *DeleteStorageTask) run() {
	log.Debugf("Task <%s> for storage <%s> started",
		"DeleteStorage", t.GetStorageName())

	if t.GetForce() {
		store, err := t.GetService().Get(t.GetStorageName())
		if err != nil {
			t.TriggerFault(types.NewErrUnhandled(err))
			return
		}

		deleteDir := NewDeleteDir(t)
		deleteDir.SetPath("")
		deleteDir.SetStorage(store)

		t.GetScheduler().Sync(deleteDir)

		// TODO: remove all segments.
	}

	err := t.GetService().Delete(t.GetStorageName())
	if err != nil {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}

	log.Debugf("Task <%s> for storage <%s> finished",
		"DeleteStorage", t.GetStorageName())
}

func (t *DeleteSegmentTask) new() {}
func (t *DeleteSegmentTask) run() {}

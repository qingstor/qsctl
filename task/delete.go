package task

import (
	"errors"

	"github.com/Xuanwo/storage/pkg/iterator"
	typ "github.com/Xuanwo/storage/types"
	log "github.com/sirupsen/logrus"

	"github.com/yunify/qsctl/v2/pkg/types"
)

func (t *DeleteDirTask) new() {
	oc := make(chan *typ.Object)
	t.SetObjectChannel(oc)

	// set recursive for list async task to list recursively
	t.SetRecursive(true)

	t.SetPathScheduleFunc(NewDeleteFilePathParametricTask)
}

func (t *DeleteDirTask) run() {
	log.Debugf("Task <%s> for path <%s> started",
		"DeleteDir", t.GetPath())

	// TODO: check logic here

	t.GetScheduler().Sync(NewIterateFileTask(t))

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

	err := t.GetService().Delete(t.GetStorageName())
	if err != nil {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}

	log.Debugf("Task <%s> for storage <%s> finished",
		"DeleteStorage", t.GetStorageName())
}

func (t *DeleteStorageForceTask) new() {}

func (t *DeleteStorageForceTask) run() {
	log.Debugf("Task <%s> for storage <%s> started",
		"DeleteStorageForce", t.GetStorageName())

	// TODO: delete all dir, delete all segments
	t.GetScheduler().Sync(NewDeleteDirTask(t))
	t.GetScheduler().Async(NewDeleteSegmentTask(t))

	log.Debugf("Task <%s> for storage <%s> finished",
		"DeleteStorageForce", t.GetStorageName())
}

func (t *DeleteSegmentTask) new() {}

func (t *DeleteSegmentTask) run() {
	log.Debugf("Task <%s> for Bucket <%s> started.", "DeleteSegment", t.GetStorageName())
	it := t.GetStorage().ListSegments("")

	for {
		o, err := it.Next()
		if err != nil && errors.Is(err, iterator.ErrDone) {
			break
		}
		if err != nil {
			t.TriggerFault(types.NewErrUnhandled(err))
			return
		}
		if err := t.GetStorage().AbortSegment(o.ID); err != nil {
			t.TriggerFault(types.NewErrUnhandled(err))
			return
		}
	}

	log.Debugf("Task <%s> for Bucket <%s> finished.", "DeleteSegment", t.GetStorageName())
}

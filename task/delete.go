package task

import (
	"github.com/yunify/qsctl/v2/pkg/types"
)

func (t *DeleteDirTask) new() {}

func (t *DeleteDirTask) run() {
	x := NewIterateFile(t)
	x.SetPathFunc(func(key string) {
		sf := NewDeleteFile(t)
		sf.SetPath(key)
		t.GetScheduler().Async(sf)
	})
	x.SetRecursive(true)
	t.GetScheduler().Sync(x)
}

func (t *DeleteFileTask) new() {}

func (t *DeleteFileTask) run() {
	if err := t.GetStorage().Delete(t.GetPath()); err != nil {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}
}

func (t *DeleteStorageTask) new() {
}

func (t *DeleteStorageTask) run() {
	if t.GetForce() {
		store, err := t.GetService().Get(t.GetStorageName())
		if err != nil {
			t.TriggerFault(types.NewErrUnhandled(err))
			return
		}

		deleteDir := NewDeleteDir(t)
		deleteDir.SetPath("")
		deleteDir.SetStorage(store)

		t.GetScheduler().Async(deleteDir)

		deleteSegment := NewDeleteSegmentDir(t)
		deleteSegment.SetPath("") // set path "" means delete all segments
		deleteSegment.SetStorage(store)

		t.GetScheduler().Async(deleteSegment)
		t.GetScheduler().Wait()
	}

	err := t.GetService().Delete(t.GetStorageName())
	if err != nil {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}
}

func (t *DeleteSegmentTask) new() {}
func (t *DeleteSegmentTask) run() {
	if err := t.GetStorage().AbortSegment(t.GetSegmentID()); err != nil {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}
}

func (t *DeleteSegmentDirTask) new() {}

func (t *DeleteSegmentDirTask) run() {
	x := NewIterateSegment(t)
	x.SetSegmentIDFunc(func(id string) {
		sf := NewDeleteSegment(t)
		sf.SetSegmentID(id)
		t.GetScheduler().Async(sf)
	})
	t.GetScheduler().Sync(x)
}

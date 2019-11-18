package task

import (
	"errors"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/pkg/segment"
	typ "github.com/Xuanwo/storage/types"
	"github.com/yunify/qsctl/v2/pkg/types"
)

func (t *DeleteDirTask) new() {}

func (t *DeleteDirTask) run() {
	x := NewListDir(t)
	x.SetFileFunc(func(o *typ.Object) {
		sf := NewDeleteFile(t)
		sf.SetPath(o.Name)
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
	segmenter, ok := t.GetStorage().(storage.Segmenter)
	if !ok {
		t.TriggerFault(types.NewErrUnhandled(errors.New("no supported")))
		return
	}

	if err := segmenter.AbortSegment(t.GetSegmentID()); err != nil {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}
}

func (t *DeleteSegmentDirTask) new() {}

func (t *DeleteSegmentDirTask) run() {
	x := NewListSegment(t)
	x.SetSegmentFunc(func(s *segment.Segment) {
		sf := NewDeleteSegment(t)
		sf.SetSegmentID(s.ID)
		t.GetScheduler().Async(sf)
	})
	t.GetScheduler().Sync(x)
}

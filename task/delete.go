package task

import (
	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/pkg/segment"
	typ "github.com/Xuanwo/storage/types"

	"github.com/yunify/qsctl/v2/pkg/types"
)

func (t *DeleteFileTask) new() {}
func (t *DeleteFileTask) run() {
	if err := t.GetStorage().Delete(t.GetPath()); err != nil {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}
}

func (t *DeleteDirTask) new() {}
func (t *DeleteDirTask) run() {
	x := NewListDir(t)
	x.SetFileFunc(func(o *typ.Object) {
		sf := NewDeleteFile(t)
		sf.SetPath(o.Name)
		t.GetScheduler().Async(sf)
	})
	x.SetDirFunc(func(o *typ.Object) {
		sf := NewDeleteDir(t)
		sf.SetPath(o.Name)
		t.GetScheduler().Sync(sf)
	})
	t.GetScheduler().Sync(x)
}

func (t *DeleteSegmentTask) new() {}
func (t *DeleteSegmentTask) run() {
	if err := t.GetSegmenter().AbortSegment(t.GetSegmentID()); err != nil {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}
}

func (t *DeleteStorageTask) new() {}
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

		segmenter, ok := store.(storage.Segmenter)
		if ok {
			listSegments := NewListSegment(t)
			listSegments.SetSegmenter(segmenter)
			listSegments.SetPath("")
			listSegments.SetSegmentFunc(func(s *segment.Segment) {
				sf := NewDeleteSegment(t)
				sf.SetSegmentID(s.ID)

				t.GetScheduler().Async(sf)
			})

			t.GetScheduler().Async(listSegments)
		}

		t.GetScheduler().Wait()
	}

	err := t.GetService().Delete(t.GetStorageName())
	if err != nil {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}
}

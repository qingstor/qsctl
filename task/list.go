package task

import (
	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/types/pairs"
	"github.com/yunify/qsctl/v2/pkg/types"
)

func (t *ListDirTask) new() {}

func (t *ListDirTask) run() {
	err := t.GetStorage().ListDir(
		t.GetPath(),
		pairs.WithDirFunc(t.GetDirFunc()),
		pairs.WithFileFunc(t.GetFileFunc()),
	)
	if err != nil {
		t.TriggerFault(err)
		return
	}
}

func (t *ListSegmentTask) new() {}

func (t *ListSegmentTask) run() {
	err := t.GetSegmenter().ListSegments(t.GetPath(),
		pairs.WithSegmentFunc(t.GetSegmentFunc()))
	if err != nil {
		t.TriggerFault(err)
		return
	}
}

func (t *ListStorageTask) new() {}
func (t *ListStorageTask) run() {
	err := t.GetService().List(pairs.WithLocation(t.GetZone()), pairs.WithStoragerFunc(func(storager storage.Storager) {
		t.GetStoragerFunc()(storager)
	}))
	if err != nil {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}
}

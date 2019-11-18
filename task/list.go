package task

import (
	"errors"

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
	segmenter, ok := t.GetStorage().(storage.Segmenter)
	if !ok {
		t.TriggerFault(types.NewErrUnhandled(errors.New("no supported")))
		return
	}

	err := segmenter.ListSegments(t.GetPath(),
		pairs.WithSegmentFunc(t.GetSegmentFunc()))
	if err != nil {
		t.TriggerFault(err)
		return
	}
}

func (t *ListStorageTask) new() {}
func (t *ListStorageTask) run() {
	err := t.GetService().List(pairs.WithLocation(t.GetZone()), pairs.WithStoragerFunc(func(storager storage.Storager) {
		t.GetStoragerFunc()
	}))
	if err != nil {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}
}

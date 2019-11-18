package task

import (
	"errors"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/types/pairs"
	"github.com/yunify/qsctl/v2/pkg/types"
)

func (t *ReachFileTask) new() {}
func (t *ReachFileTask) run() {
	reacher, ok := t.GetStorage().(storage.Reacher)
	if !ok {
		// TODO: we need a better error
		t.TriggerFault(types.NewErrUnhandled(errors.New("no supported")))
		return
	}
	url, err := reacher.Reach(t.GetPath(), pairs.WithExpire(t.GetExpire()))
	if err != nil {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}
	t.SetURL(url)
}

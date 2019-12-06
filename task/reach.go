package task

import (
	"github.com/Xuanwo/storage/types/pairs"

	"github.com/qingstor/qsctl/v2/pkg/types"
)

func (t *ReachFileTask) new() {}
func (t *ReachFileTask) run() {
	url, err := t.GetReacher().Reach(t.GetPath(), pairs.WithExpire(t.GetExpire()))
	if err != nil {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}
	t.SetURL(url)
}

package task

import (
	typ "github.com/Xuanwo/storage/types"
	"github.com/yunify/qsctl/v2/pkg/types"
)

func (t *ReachFileTask) new() {}
func (t *ReachFileTask) run() {
	url, err := t.GetStorage().Reach(t.GetPath(), typ.WithExpire(t.GetExpire()))
	if err != nil {
		t.TriggerFault(types.NewErrUnhandled(err))
		return
	}
	t.SetURL(url)
}

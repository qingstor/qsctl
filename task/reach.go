package task

import (
	typ "github.com/Xuanwo/storage/types"
	log "github.com/sirupsen/logrus"

	"github.com/yunify/qsctl/v2/pkg/fault"
)

func (t *ReachTask) new() {}
func (t *ReachTask) run() {
	url, err := t.GetStorage().Reach(t.GetPath(), typ.WithExpire(t.GetExpire()))
	if err != nil {
		t.TriggerFault(fault.NewUnhandled(err))
		return
	}
	t.SetURL(url)
	log.Debugf("Task <%s> for key <%s> finished, get signed URL <%s>",
		"ObjectPresignTask", t.GetPath(), "")
}

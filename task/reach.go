package task

import (
	typ "github.com/Xuanwo/storage/types"
	log "github.com/sirupsen/logrus"
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
	log.Debugf("Task <%s> for key <%s> finished, get signed URL <%s>",
		"ObjectPresignTask", t.GetPath(), "")
}

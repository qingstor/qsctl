package common

import (
	"github.com/Xuanwo/storage/types"
	log "github.com/sirupsen/logrus"

	"github.com/yunify/qsctl/v2/pkg/fault"
)

func (t *ObjectPresignTask) run() {
	url, err := t.GetDestinationStorage().Reach(t.GetKey(), types.WithExpire(t.GetExpire()))
	if err != nil {
		t.TriggerFault(fault.NewUnhandled(err))
		return
	}
	t.SetURL(url)
	log.Debugf("Task <%s> for key <%s> finished, get signed URL <%s>",
		"ObjectPresignTask", t.GetKey(), "")
}

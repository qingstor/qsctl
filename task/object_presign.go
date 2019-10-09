package task

import (
	log "github.com/sirupsen/logrus"

	"github.com/yunify/qsctl/v2/pkg/fault"
)

func (t *ObjectPresignTask) run() {
	if _, err := t.GetStorage().HeadObject(t.GetKey()); err != nil {
		t.TriggerFault(fault.NewUnhandled(err))
		return
	}

	url, err := t.GetStorage().PresignObject(t.GetKey(), t.GetExpire())
	if err != nil {
		t.TriggerFault(fault.NewUnhandled(err))
		return
	}
	t.SetURL(url)
	log.Debugf("Task <%s> for key <%s> finished, get signed URL <%s>",
		"ObjectPresignTask", t.GetKey(), url)
}

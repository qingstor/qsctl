package common

import (
	log "github.com/sirupsen/logrus"
	"github.com/yunify/qsctl/v2/pkg/fault"
)

func (t *ObjectPresignTask) run() {
	// TODO: add expire support.
	url, err := t.GetDestinationStorage().Reach(t.GetKey())
	if err != nil {
		t.TriggerFault(fault.NewUnhandled(err))
		return
	}
	t.SetURL(url)
	log.Debugf("Task <%s> for key <%s> finished, get signed URL <%s>",
		"ObjectPresignTask", t.GetKey(), "")
}

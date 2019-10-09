package task

import (
	log "github.com/sirupsen/logrus"

	"github.com/yunify/qsctl/v2/pkg/fault"
)

func (t *ObjectDeleteTask) run() {
	if _, err := t.GetStorage().HeadObject(t.GetKey()); err != nil {
		t.TriggerFault(fault.NewUnhandled(err))
		return
	}

	if err := t.GetStorage().DeleteObject(t.GetKey()); err != nil {
		t.TriggerFault(fault.NewUnhandled(err))
		return
	}

	log.Debugf("Task <%s> for key <%s> finished",
		"ObjectDeleteTask", t.GetKey())
}

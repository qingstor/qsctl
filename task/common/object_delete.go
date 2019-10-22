package common

import (
	log "github.com/sirupsen/logrus"

	"github.com/yunify/qsctl/v2/pkg/fault"
)

func (t *ObjectDeleteTask) run() {
	if err := t.GetDestinationStorage().Delete(t.GetDestinationPath()); err != nil {
		t.TriggerFault(fault.NewUnhandled(err))
		return
	}

	log.Debugf("Task <%s> for key <%s> finished",
		"ObjectDeleteTask", t.GetDestinationPath())
}

package common

import (
	log "github.com/sirupsen/logrus"

	"github.com/yunify/qsctl/v2/pkg/fault"
)

func (t *ObjectStatTask) run() {
	om, err := t.GetDestinationStorage().Stat(t.GetDestinationPath())
	if err != nil {
		t.TriggerFault(fault.NewUnhandled(err))
		return
	}
	t.SetObject(om)
	log.Debugf("Task <%s> for Key <%s> finished.", "StatObjectTask", t.GetDestinationPath())
}

package task

import (
	log "github.com/sirupsen/logrus"
	"github.com/yunify/qsctl/v2/pkg/fault"
)

func (t *StatTask) new() {}
func (t *StatTask) run() {
	om, err := t.GetStorage().Stat(t.GetPath())
	if err != nil {
		t.TriggerFault(fault.NewUnhandled(err))
		return
	}
	t.SetObject(om)
	log.Debugf("Task <%s> for Key <%s> finished.", "StatObjectTask", t.GetPath())
}

package task

import (
	log "github.com/sirupsen/logrus"
	"github.com/yunify/qsctl/v2/pkg/fault"
)

func (t *ObjectStatTask) run() {
	om, err := t.GetStorage().HeadObject(t.GetKey())
	if err != nil {
		t.TriggerError(fault.NewUnhandled(err))
		return
	}
	// replace the original om
	t.SetObjectMeta(om)
	log.Debugf("Task <%s> for Key <%s> finished.", "StatObjectTask", t.GetKey())
}

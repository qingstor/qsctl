package task

import (
	"github.com/Xuanwo/navvy"
	"github.com/yunify/qsctl/v2/pkg/fault"
)

// NewStatTask will create a stat task.
func NewStatTask(fn func(*StatTask)) *StatTask {
	t := &StatTask{}

	pool := navvy.NewPool(10)
	t.SetPool(pool)

	fn(t)
	t.AddTODOs(NewObjectStatTask)
	return t
}
func (t *ObjectStatTask) run() {
	om, err := t.GetDestinationStorage().Stat(t.GetDestinationPath())
	if err != nil {
		t.TriggerFault(fault.NewUnhandled(err))
		return
	}
	t.SetObject(om)
	log.Debugf("Task <%s> for Key <%s> finished.", "StatObjectTask", t.GetDestinationPath())
}

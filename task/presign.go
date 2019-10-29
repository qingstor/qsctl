package task

import (
	"github.com/Xuanwo/navvy"
	"github.com/yunify/qsctl/v2/pkg/fault"
)

// NewPresignTask will create a presign task
func NewPresignTask(fn func(*PresignTask)) *PresignTask {
	t := &PresignTask{}

	pool := navvy.NewPool(10)
	t.SetPool(pool)

	fn(t)

	t.AddTODOs(NewObjectPresignTask)
	return t
}
func (t *ObjectPresignTask) run() {
	url, err := t.GetDestinationStorage().Reach(t.GetDestinationPath(), typ.WithExpire(t.GetExpire()))
	if err != nil {
		t.TriggerFault(fault.NewUnhandled(err))
		return
	}
	t.SetURL(url)
	log.Debugf("Task <%s> for key <%s> finished, get signed URL <%s>",
		"ObjectPresignTask", t.GetDestinationPath(), "")
}

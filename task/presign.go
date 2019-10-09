package task

import (
	"github.com/Xuanwo/navvy"

	"github.com/yunify/qsctl/v2/pkg/fault"
)

// NewPresignTask will create a presign task
func NewPresignTask(fn func(*PresignTask)) *PresignTask {
	t := &PresignTask{}

	pool, err := navvy.NewPool(10)
	if err != nil {
		t.TriggerFault(fault.NewUnhandled(err))
		return t
	}
	t.SetPool(pool)

	fn(t)

	t.AddTODOs(NewObjectPresignTask)
	return t
}

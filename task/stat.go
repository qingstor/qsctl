package task

import (
	"github.com/Xuanwo/navvy"
	"github.com/yunify/qsctl/v2/pkg/fault"
)

// NewStatTask will create a stat task.
func NewStatTask(fn func(*StatTask)) *StatTask {
	t := &StatTask{}

	pool, err := navvy.NewPool(10)
	if err != nil {
		t.TriggerError(fault.NewUnhandled(err))
		return t
	}
	t.SetPool(pool)

	fn(t)
	t.AddTODOs(NewObjectStatTask)
	return t
}

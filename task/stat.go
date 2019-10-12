package task

import (
	"github.com/Xuanwo/navvy"

	"github.com/yunify/qsctl/v2/task/common"
)

// NewStatTask will create a stat task.
func NewStatTask(fn func(*StatTask)) *StatTask {
	t := &StatTask{}

	pool := navvy.NewPool(10)
	t.SetPool(pool)

	fn(t)
	t.AddTODOs(common.NewObjectStatTask)
	return t
}

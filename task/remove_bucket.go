package task

import (
	"github.com/Xuanwo/navvy"

	"github.com/yunify/qsctl/v2/task/common"
)

// NewRemoveBucketTask will create a remove bucket task
func NewRemoveBucketTask(fn func(*RemoveBucketTask)) *RemoveBucketTask {
	t := &RemoveBucketTask{}

	pool := navvy.NewPool(10)
	t.SetPool(pool)

	fn(t)

	if t.ValidateFault() {
		return t
	}

	// TODO: check force flag, if true, do rm -r, then delete bucket
	t.AddTODOs(common.NewBucketDeleteTask)
	return t
}

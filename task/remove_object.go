package task

import (
	"github.com/Xuanwo/navvy"

	"github.com/yunify/qsctl/v2/task/common"
)

// NewRemoveObjectTask will create a remove object task
func NewRemoveObjectTask(fn func(*RemoveObjectTask)) *RemoveObjectTask {
	t := &RemoveObjectTask{}

	pool := navvy.NewPool(10)
	t.SetPool(pool)

	fn(t)

	if t.ValidateFault() {
		return t
	}

	// TODO: check recursive flag, if true, do rm -r
	if t.GetRecursive() {
		t.SetDeleteKey(t.GetKey())
		t.AddTODOs(common.NewRemoveDirTask)
		return t
	}
	t.AddTODOs(common.NewObjectDeleteTask)
	return t
}

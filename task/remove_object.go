package task

import (
	"github.com/Xuanwo/navvy"
)

// NewRemoveObjectTask will create a remove object task
func NewRemoveObjectTask(fn func(*RemoveObjectTask)) *RemoveObjectTask {
	t := &RemoveObjectTask{}

	pool, _ := navvy.NewPool(10)
	t.SetPool(pool)

	fn(t)

	if t.ValidateFault() {
		return t
	}

	// TODO: check recursive flag, if true, do rm -r
	if t.GetRecursive() {
		panic("rm -r not implemented.")
	}
	t.AddTODOs(NewObjectDeleteTask)
	return t
}

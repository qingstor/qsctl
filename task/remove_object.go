package task

import (
	"github.com/Xuanwo/navvy"
	storType "github.com/Xuanwo/storage/types"

	"github.com/yunify/qsctl/v2/pkg/types"
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

	// check recursive flag, if true, do rm -r
	if t.GetRecursive() {
		t.SetPrefix(t.GetKey())
		t.AddTODOs(NewRemoveDirTask)
		return t
	}
	t.AddTODOs(common.NewObjectDeleteTask)
	return t
}

func (t *RemoveDirTask) new() {
	oc := make(chan *storType.Object, 1)
	t.SetObjectChannel(oc)

	t.SetScheduler(types.NewScheduler(common.NewObjectDeleteRecursivelyTask))

	t.AddTODOs(
		common.NewObjectInitDirDeleteTask,
		common.NewObjectListTask,
		common.NewWaitTask,
	)
}

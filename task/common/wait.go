package common

import (
	"github.com/Xuanwo/navvy"
	"github.com/yunify/qsctl/v2/task/utils"

	"github.com/yunify/qsctl/v2/task/types"
)

// WaitTaskRequirement is the requirement for execute WaitTask.
type WaitTaskRequirement interface {
	navvy.Task
	types.Todoist

	types.WaitGroupGetter
	types.PoolGetter
}

// WaitTask will execute Wait Task.
type WaitTask struct {
	WaitTaskRequirement
}

// NewWaitTask will create a new Task.
func NewWaitTask(t types.Todoist) navvy.Task {
	o, ok := t.(WaitTaskRequirement)
	if !ok {
		panic("task is not fill NewWaitTask")
	}

	return &WaitTask{o}
}

// Run implement navvy.Task.
func (t *WaitTask) Run() {
	t.GetWaitGroup().Wait()

	go utils.SubmitNextTask(t.WaitTaskRequirement)
}

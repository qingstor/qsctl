// Code generated by go generate; DO NOT EDIT.
package common

import (
	"fmt"

	"github.com/Xuanwo/navvy"
	"github.com/google/uuid"

	"github.com/yunify/qsctl/v2/pkg/types"
	"github.com/yunify/qsctl/v2/utils"
)

var _ navvy.Pool
var _ types.Pool
var _ = utils.SubmitNextTask
var _ = uuid.New()

// waitTaskRequirement is the requirement for execute WaitTask.
type waitTaskRequirement interface {
	navvy.Task
	types.Todoist
	types.PoolGetter
	types.FaultSetter
	types.FaultValidator
	types.IDGetter

	// Inherited value
	types.SchedulerGetter
	// Runtime value
}

// mockWaitTask is the mock task for WaitTask.
type mockWaitTask struct {
	types.Todo
	types.Pool
	types.Fault
	types.ID

	// Inherited value
	types.Scheduler
	// Runtime value
}

func (t *mockWaitTask) Run() {
	panic("mockWaitTask should not be run.")
}

// WaitTask will wait until parent task finished.
type WaitTask struct {
	waitTaskRequirement
}

// Run implement navvy.Task.
func (t *WaitTask) Run() {
	t.run()
	if t.ValidateFault() {
		return
	}
	utils.SubmitNextTask(t.waitTaskRequirement)
}

func (t *WaitTask) TriggerFault(err error) {
	t.SetFault(fmt.Errorf("Task Wait failed: {%w}", err))
}

// NewWaitTask will create a new WaitTask.
func NewWaitTask(task types.Todoist) navvy.Task {
	return &WaitTask{task.(waitTaskRequirement)}
}
// Code generated by go generate; DO NOT EDIT.
package common

import (
	"fmt"

	"github.com/Xuanwo/navvy"
	"github.com/google/uuid"

	"github.com/yunify/qsctl/v2/pkg/types"
)

var _ navvy.Pool
var _ types.Pool
var _ = uuid.New()

// doneSchedulerTaskRequirement is the requirement for execute DoneSchedulerTask.
type doneSchedulerTaskRequirement interface {
	navvy.Task

	// Inherited value

	// Mutable value
}

// mockDoneSchedulerTask is the mock task for DoneSchedulerTask.
type mockDoneSchedulerTask struct {
	types.Pool
	types.Fault
	types.ID

	// Inherited value

	// Mutable value
}

func (t *mockDoneSchedulerTask) Run() {
	panic("mockDoneSchedulerTask should not be run.")
}

// DoneSchedulerTask will done scheduler manually.
type DoneSchedulerTask struct {
	doneSchedulerTaskRequirement

	// Predefined runtime value
	types.Fault
	types.ID
	types.Scheduler

	// Runtime value
}

// Run implement navvy.Task
func (t *DoneSchedulerTask) Run() {
	t.run()
}

func (t *DoneSchedulerTask) TriggerFault(err error) {
	t.SetFault(fmt.Errorf("Task DoneScheduler failed: {%w}", err))
}

// NewDoneSchedulerTask will create a DoneSchedulerTask and fetch inherited data from parent task.
func NewDoneSchedulerTask(task navvy.Task) navvy.Task {
	t := &DoneSchedulerTask{
		doneSchedulerTaskRequirement: task.(doneSchedulerTaskRequirement),
	}
	t.SetID(uuid.New().String())
	t.new()
	return t
}

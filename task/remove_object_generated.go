// Code generated by go generate; DO NOT EDIT.
package task

import (
	"fmt"

	"github.com/Xuanwo/navvy"
	"github.com/google/uuid"

	"github.com/yunify/qsctl/v2/pkg/types"
)

var _ navvy.Pool
var _ types.Pool
var _ = uuid.New()

// removeObjectTaskRequirement is the requirement for execute RemoveObjectTask.
type removeObjectTaskRequirement interface {
	navvy.Task

	// Inherited value

	// Mutable value
}

// mockRemoveObjectTask is the mock task for RemoveObjectTask.
type mockRemoveObjectTask struct {
	types.Pool
	types.Fault
	types.ID

	// Inherited value

	// Mutable value
}

func (t *mockRemoveObjectTask) Run() {
	panic("mockRemoveObjectTask should not be run.")
}

// RemoveObjectTask will will remove object.
type RemoveObjectTask struct {
	removeObjectTaskRequirement

	// Predefined runtime value
	types.Fault
	types.ID
	types.Scheduler

	// Runtime value
	types.DestinationPath
	types.DestinationStorage
	types.DestinationType
	types.Pool
	types.Recursive
}

// Run implement navvy.Task
func (t *RemoveObjectTask) Run() {
	t.run()
}

func (t *RemoveObjectTask) TriggerFault(err error) {
	t.SetFault(fmt.Errorf("Task RemoveObject failed: {%w}", err))
}

// NewRemoveObjectTask will create a RemoveObjectTask and fetch inherited data from parent task.
func NewRemoveObjectTask(task navvy.Task) navvy.Task {
	t := &RemoveObjectTask{
		removeObjectTaskRequirement: task.(removeObjectTaskRequirement),
	}
	t.SetID(uuid.New().String())
	t.new()
	return t
}

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

// statTaskRequirement is the requirement for execute StatTask.
type statTaskRequirement interface {
	navvy.Task

	// Inherited value

	// Mutable value
}

// mockStatTask is the mock task for StatTask.
type mockStatTask struct {
	types.Pool
	types.Fault
	types.ID

	// Inherited value

	// Mutable value
}

func (t *mockStatTask) Run() {
	panic("mockStatTask should not be run.")
}

// StatTask will will stat a remote object.
type StatTask struct {
	statTaskRequirement

	// Predefined runtime value
	types.Fault
	types.ID
	types.Scheduler

	// Runtime value
	types.DestinationPath
	types.DestinationStorage
	types.DestinationType
	types.Object
	types.Pool
}

// Run implement navvy.Task
func (t *StatTask) Run() {
	t.run()
}

func (t *StatTask) TriggerFault(err error) {
	t.SetFault(fmt.Errorf("Task Stat failed: {%w}", err))
}

// NewStatTask will create a StatTask and fetch inherited data from parent task.
func NewStatTask(task navvy.Task) navvy.Task {
	t := &StatTask{
		statTaskRequirement: task.(statTaskRequirement),
	}
	t.SetID(uuid.New().String())
	t.new()
	return t
}

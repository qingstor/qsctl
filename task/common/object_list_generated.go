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

// objectListTaskRequirement is the requirement for execute ObjectListTask.
type objectListTaskRequirement interface {
	navvy.Task

	// Inherited value
	types.DestinationPathGetter
	types.DestinationStorageGetter
	types.ObjectChannelGetter
	types.RecursiveGetter

	// Mutable value
}

// mockObjectListTask is the mock task for ObjectListTask.
type mockObjectListTask struct {
	types.Todo
	types.Pool
	types.Fault
	types.ID

	// Inherited value
	types.DestinationPath
	types.DestinationStorage
	types.ObjectChannel
	types.Recursive

	// Mutable value
}

func (t *mockObjectListTask) Run() {
	panic("mockObjectListTask should not be run.")
}

// ObjectListTask will list objects.
type ObjectListTask struct {
	objectListTaskRequirement

	// Predefined runtime value
	types.Fault
	types.ID
	types.Scheduler

	// Runtime value
}

// Run implement navvy.Task
func (t *ObjectListTask) Run() {
	t.run()
}

func (t *ObjectListTask) TriggerFault(err error) {
	t.SetFault(fmt.Errorf("Task ObjectList failed: {%w}", err))
}

// Wait will wait until ObjectListTask has been finished
func (t *ObjectListTask) Wait() {
	t.GetPool().Wait()
}

// objectListAsyncTaskRequirement is the requirement for execute ObjectListAsyncTask.
type objectListAsyncTaskRequirement interface {
	navvy.Task

	// Inherited value
	types.DestinationPathGetter
	types.DestinationStorageGetter
	types.ObjectChannelGetter
	types.RecursiveGetter

	// Mutable value
}

// mockObjectListAsyncTask is the mock task for ObjectListAsyncTask.
type mockObjectListAsyncTask struct {
	types.Todo
	types.Pool
	types.Fault
	types.ID

	// Inherited value
	types.DestinationPath
	types.DestinationStorage
	types.ObjectChannel
	types.Recursive

	// Mutable value
}

func (t *mockObjectListAsyncTask) Run() {
	panic("mockObjectListAsyncTask should not be run.")
}

// ObjectListAsyncTask will list objects.
type ObjectListAsyncTask struct {
	objectListAsyncTaskRequirement

	// Predefined runtime value
	types.Fault
	types.ID
	types.Scheduler

	// Runtime value
}

// Run implement navvy.Task
func (t *ObjectListAsyncTask) Run() {
	t.run()
}

func (t *ObjectListAsyncTask) TriggerFault(err error) {
	t.SetFault(fmt.Errorf("Task ObjectListAsync failed: {%w}", err))
}

// Wait will wait until ObjectListAsyncTask has been finished
func (t *ObjectListAsyncTask) Wait() {
	t.GetPool().Wait()
}

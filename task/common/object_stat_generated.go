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

// objectStatTaskRequirement is the requirement for execute ObjectStatTask.
type objectStatTaskRequirement interface {
	navvy.Task
	types.Todoist
	types.PoolGetter
	types.FaultSetter
	types.FaultValidator
	types.IDGetter

	// Inherited value
	types.DestinationStorageGetter
	types.KeyGetter
	// Runtime value
	types.ObjectSetter
}

// mockObjectStatTask is the mock task for ObjectStatTask.
type mockObjectStatTask struct {
	types.Todo
	types.Pool
	types.Fault
	types.ID

	// Inherited value
	types.DestinationStorage
	types.Key
	// Runtime value
	types.Object
}

func (t *mockObjectStatTask) Run() {
	panic("mockObjectStatTask should not be run.")
}

// ObjectStatTask will stat a remote object by request headObject.
type ObjectStatTask struct {
	objectStatTaskRequirement
}

// Run implement navvy.Task.
func (t *ObjectStatTask) Run() {
	t.run()
	if t.ValidateFault() {
		return
	}
	utils.SubmitNextTask(t.objectStatTaskRequirement)
}

func (t *ObjectStatTask) TriggerFault(err error) {
	t.SetFault(fmt.Errorf("Task ObjectStat failed: {%w}", err))
}

// NewObjectStatTask will create a new ObjectStatTask.
func NewObjectStatTask(task types.Todoist) navvy.Task {
	return &ObjectStatTask{task.(objectStatTaskRequirement)}
}

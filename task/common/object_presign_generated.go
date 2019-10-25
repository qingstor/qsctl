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

// objectPresignTaskRequirement is the requirement for execute ObjectPresignTask.
type objectPresignTaskRequirement interface {
	navvy.Task
	types.Todoist
	types.PoolGetter
	types.FaultSetter
	types.FaultValidator
	types.IDGetter

	// Inherited value
	types.BucketNameGetter
	types.DestinationPathGetter
	types.DestinationStorageGetter
	types.ExpireGetter
	// Runtime value
	types.URLSetter
}

// mockObjectPresignTask is the mock task for ObjectPresignTask.
type mockObjectPresignTask struct {
	types.Todo
	types.Pool
	types.Fault
	types.ID

	// Inherited value
	types.BucketName
	types.DestinationPath
	types.DestinationStorage
	types.Expire
	// Runtime value
	types.URL
}

func (t *mockObjectPresignTask) Run() {
	panic("mockObjectPresignTask should not be run.")
}

// ObjectPresignTask will will presign a remote object and return the signed url.
type ObjectPresignTask struct {
	objectPresignTaskRequirement
}

// Run implement navvy.Task.
func (t *ObjectPresignTask) Run() {
	t.run()
	if t.ValidateFault() {
		return
	}
	utils.SubmitNextTask(t.objectPresignTaskRequirement)
}

func (t *ObjectPresignTask) TriggerFault(err error) {
	t.SetFault(fmt.Errorf("Task ObjectPresign failed: {%w}", err))
}

// NewObjectPresignTask will create a new ObjectPresignTask.
func NewObjectPresignTask(task types.Todoist) navvy.Task {
	return &ObjectPresignTask{task.(objectPresignTaskRequirement)}
}
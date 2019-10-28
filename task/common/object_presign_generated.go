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

// objectPresignTaskRequirement is the requirement for execute ObjectPresignTask.
type objectPresignTaskRequirement interface {
	navvy.Task

	// Inherited value
	types.BucketNameGetter
	types.DestinationPathGetter
	types.DestinationStorageGetter
	types.ExpireGetter

	// Mutable value
}

// mockObjectPresignTask is the mock task for ObjectPresignTask.
type mockObjectPresignTask struct {
	types.Pool
	types.Fault
	types.ID

	// Inherited value
	types.BucketName
	types.DestinationPath
	types.DestinationStorage
	types.Expire

	// Mutable value
}

func (t *mockObjectPresignTask) Run() {
	panic("mockObjectPresignTask should not be run.")
}

// ObjectPresignTask will will presign a remote object and return the signed url.
type ObjectPresignTask struct {
	objectPresignTaskRequirement

	// Predefined runtime value
	types.Fault
	types.ID
	types.Scheduler

	// Runtime value
	types.URL
}

// Run implement navvy.Task
func (t *ObjectPresignTask) Run() {
	t.run()
}

func (t *ObjectPresignTask) TriggerFault(err error) {
	t.SetFault(fmt.Errorf("Task ObjectPresign failed: {%w}", err))
}

// NewObjectPresignTask will create a ObjectPresignTask and fetch inherited data from Task.
func NewObjectPresignTask(task navvy.Task) navvy.Task {
	t := &ObjectPresignTask{
		objectPresignTaskRequirement: task.(objectPresignTaskRequirement),
	}
	t.SetID(uuid.New().String())
	t.new()
	return t
}

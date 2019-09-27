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

// bucketCreateTaskRequirement is the requirement for execute BucketCreateTask.
type bucketCreateTaskRequirement interface {
	navvy.Task
	types.Todoist
	types.PoolGetter
	types.FaultSetter
	types.FaultValidator
	types.IDGetter

	// Inherited value
	types.BucketNameGetter
	types.StorageGetter
	// Runtime value
}

// mockBucketCreateTask is the mock task for BucketCreateTask.
type mockBucketCreateTask struct {
	types.Todo
	types.Pool
	types.Fault
	types.ID

	// Inherited value
	types.BucketName
	types.Storage
	// Runtime value
}

func (t *mockBucketCreateTask) Run() {
	panic("mockBucketCreateTask should not be run.")
}

// BucketCreateTask will send put request to create a bucket.
type BucketCreateTask struct {
	bucketCreateTaskRequirement
}

// Run implement navvy.Task.
func (t *BucketCreateTask) Run() {
	t.run()
	if t.ValidateFault() {
		return
	}
	utils.SubmitNextTask(t.bucketCreateTaskRequirement)
}

func (t *BucketCreateTask) TriggerFault(err error) {
	t.SetFault(fmt.Errorf("Task BucketCreate failed: {%w}", err))
}

// NewBucketCreateTask will create a new BucketCreateTask.
func NewBucketCreateTask(task types.Todoist) navvy.Task {
	return &BucketCreateTask{task.(bucketCreateTaskRequirement)}
}

// Code generated by go generate; DO NOT EDIT.
package common

import (
	"github.com/Xuanwo/navvy"

	"github.com/yunify/qsctl/v2/pkg/types"
	"github.com/yunify/qsctl/v2/task/utils"
)

var _ navvy.Pool
var _ types.Pool
var _ = utils.SubmitNextTask

// bucketCreateTaskRequirement is the requirement for execute BucketCreateTask.
type bucketCreateTaskRequirement interface {
	navvy.Task
	types.Todoist
	types.PoolGetter

	// Inherited value
	types.BucketNameGetter
	types.StorageGetter
	// Runtime value
}

// mockBucketCreateTask is the mock task for BucketCreateTask.
type mockBucketCreateTask struct {
	types.Todo
	types.Pool

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
	utils.SubmitNextTask(t.bucketCreateTaskRequirement)
}

// NewBucketCreateTask will create a new BucketCreateTask.
func NewBucketCreateTask(task types.Todoist) navvy.Task {
	return &BucketCreateTask{task.(bucketCreateTaskRequirement)}
}

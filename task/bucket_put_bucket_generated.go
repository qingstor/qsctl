// Code generated by go generate; DO NOT EDIT.
package task

import (
	"github.com/Xuanwo/navvy"

	"github.com/yunify/qsctl/v2/task/types"
	"github.com/yunify/qsctl/v2/task/utils"
)

var _ navvy.Pool
var _ types.Pool
var _ = utils.SubmitNextTask

// putBucketTaskRequirement is the requirement for execute PutBucketTask.
type putBucketTaskRequirement interface {
	navvy.Task
	types.Todoist
	types.PoolGetter

	// Inherited value
	types.BucketNameGetter
	types.StorageGetter
	// Runtime value
}

// mockPutBucketTask is the mock task for PutBucketTask.
type mockPutBucketTask struct {
	types.Todo
	types.Pool

	// Inherited value
	types.BucketName
	types.Storage
	// Runtime value
}

func (t *mockPutBucketTask) Run() {
	panic("mockPutBucketTask should not be run.")
}

// PutBucketTask will send put request to create a bucket.
type PutBucketTask struct {
	putBucketTaskRequirement
}

// Run implement navvy.Task.
func (t *PutBucketTask) Run() {
	t.run()
	utils.SubmitNextTask(t.putBucketTaskRequirement)
}

// NewPutBucketTask will create a new PutBucketTask.
func NewPutBucketTask(task types.Todoist) navvy.Task {
	return &PutBucketTask{task.(putBucketTaskRequirement)}
}

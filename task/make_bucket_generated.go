// Code generated by go generate; DO NOT EDIT.
package task

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

// makeBucketTaskRequirement is the requirement for execute MakeBucketTask.
type makeBucketTaskRequirement interface {
	navvy.Task

	// Inherited value
}

// mockMakeBucketTask is the mock task for MakeBucketTask.
type mockMakeBucketTask struct {
	types.Todo
	types.Pool
	types.Fault
	types.ID

	// Inherited value
}

func (t *mockMakeBucketTask) Run() {
	panic("mockMakeBucketTask should not be run.")
}

// MakeBucketTask will make new bucket with given key.
type MakeBucketTask struct {
	makeBucketTaskRequirement

	// Predefined runtime value
	types.Fault
	types.ID
	types.Todo

	// Runtime value
	types.BucketName
	types.Pool
	types.Storage
	types.Zone
}

// Run implement navvy.Task
func (t *MakeBucketTask) Run() {
	if t.ValidateFault() {
		return
	}
	utils.SubmitNextTask(t)
}

func (t *MakeBucketTask) TriggerError(err error) {
	t.SetFault(fmt.Errorf("Task MakeBucket failed: {%w}", err))
}

// Wait will wait until MakeBucketTask has been finished
func (t *MakeBucketTask) Wait() {
	t.GetPool().Wait()
}

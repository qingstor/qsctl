// Code generated by go generate; DO NOT EDIT.
package task

import (
	"github.com/Xuanwo/navvy"

	"github.com/yunify/qsctl/v2/pkg/types"
	"github.com/yunify/qsctl/v2/utils"
)

var _ navvy.Pool
var _ types.Pool
var _ = utils.SubmitNextTask

// objectPresignTaskRequirement is the requirement for execute ObjectPresignTask.
type objectPresignTaskRequirement interface {
	navvy.Task
	types.Todoist
	types.PoolGetter

	// Inherited value
	types.BucketNameGetter
	types.ExpireGetter
	types.KeyGetter
	types.StorageGetter
	// Runtime value
	types.URLSetter
}

// mockObjectPresignTask is the mock task for ObjectPresignTask.
type mockObjectPresignTask struct {
	types.Todo
	types.Pool

	// Inherited value
	types.BucketName
	types.Expire
	types.Key
	types.Storage
	// Runtime value
	types.URL
}

func (t *mockObjectPresignTask) Run() {
	panic("mockObjectPresignTask should not be run.")
}

// ObjectPresignTask will will presign a remote object and return the signed url, is the sub task for Presign.
type ObjectPresignTask struct {
	objectPresignTaskRequirement
}

// Run implement navvy.Task.
func (t *ObjectPresignTask) Run() {
	t.run()
	utils.SubmitNextTask(t.objectPresignTaskRequirement)
}

// NewObjectPresignTask will create a new ObjectPresignTask.
func NewObjectPresignTask(task types.Todoist) navvy.Task {
	return &ObjectPresignTask{task.(objectPresignTaskRequirement)}
}

// objectPresignPrivateTaskRequirement is the requirement for execute ObjectPresignPrivateTask.
type objectPresignPrivateTaskRequirement interface {
	navvy.Task
	types.Todoist
	types.PoolGetter

	// Inherited value
	types.ExpireGetter
	types.KeyGetter
	types.StorageGetter
	// Runtime value
	types.URLSetter
}

// mockObjectPresignPrivateTask is the mock task for ObjectPresignPrivateTask.
type mockObjectPresignPrivateTask struct {
	types.Todo
	types.Pool

	// Inherited value
	types.Expire
	types.Key
	types.Storage
	// Runtime value
	types.URL
}

func (t *mockObjectPresignPrivateTask) Run() {
	panic("mockObjectPresignPrivateTask should not be run.")
}

// ObjectPresignPrivateTask will will presign a remote object at private bucket.
type ObjectPresignPrivateTask struct {
	objectPresignPrivateTaskRequirement
}

// Run implement navvy.Task.
func (t *ObjectPresignPrivateTask) Run() {
	t.run()
	utils.SubmitNextTask(t.objectPresignPrivateTaskRequirement)
}

// NewObjectPresignPrivateTask will create a new ObjectPresignPrivateTask.
func NewObjectPresignPrivateTask(task types.Todoist) navvy.Task {
	return &ObjectPresignPrivateTask{task.(objectPresignPrivateTaskRequirement)}
}

// objectPresignPublicTaskRequirement is the requirement for execute ObjectPresignPublicTask.
type objectPresignPublicTaskRequirement interface {
	navvy.Task
	types.Todoist
	types.PoolGetter

	// Inherited value
	types.BucketNameGetter
	types.KeyGetter
	types.StorageGetter
	// Runtime value
	types.URLSetter
}

// mockObjectPresignPublicTask is the mock task for ObjectPresignPublicTask.
type mockObjectPresignPublicTask struct {
	types.Todo
	types.Pool

	// Inherited value
	types.BucketName
	types.Key
	types.Storage
	// Runtime value
	types.URL
}

func (t *mockObjectPresignPublicTask) Run() {
	panic("mockObjectPresignPublicTask should not be run.")
}

// ObjectPresignPublicTask will will presign a remote object at public bucket.
type ObjectPresignPublicTask struct {
	objectPresignPublicTaskRequirement
}

// Run implement navvy.Task.
func (t *ObjectPresignPublicTask) Run() {
	t.run()
	utils.SubmitNextTask(t.objectPresignPublicTaskRequirement)
}

// NewObjectPresignPublicTask will create a new ObjectPresignPublicTask.
func NewObjectPresignPublicTask(task types.Todoist) navvy.Task {
	return &ObjectPresignPublicTask{task.(objectPresignPublicTaskRequirement)}
}

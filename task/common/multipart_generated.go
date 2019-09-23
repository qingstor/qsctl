// Code generated by go generate; DO NOT EDIT.
package common

import (
	"github.com/Xuanwo/navvy"

	"github.com/yunify/qsctl/v2/task/types"
	"github.com/yunify/qsctl/v2/task/utils"
)

var _ navvy.Pool
var _ types.Pool
var _ = utils.SubmitNextTask

// MultipartCompleteTaskRequirement is the requirement for execute MultipartCompleteTask.
type MultipartCompleteTaskRequirement interface {
	navvy.Task
	types.Todoist
	types.PoolGetter

	// Inherited value
	types.CurrentPartNumberGetter
	types.KeyGetter
	types.StorageGetter
	types.UploadIDGetter

	// Runtime value
}

// MultipartCompleteTask will upload a multipart via stream.
type MultipartCompleteTask struct {
	MultipartCompleteTaskRequirement
}

// mockMultipartCompleteTask is the mock task for MultipartCompleteTask.
type mockMultipartCompleteTask struct {
	types.Todo
	types.Pool

	// Inherited value
	types.CurrentPartNumber
	types.Key
	types.Storage
	types.UploadID

	// Runtime value
}

func (t *mockMultipartCompleteTask) Run() {
	panic("mockMultipartCompleteTask should not be run.")
}

// NewMultipartCompleteTask will create a new MultipartCompleteTask.
func NewMultipartCompleteTask(task types.Todoist) navvy.Task {
	return &MultipartCompleteTask{task.(MultipartCompleteTaskRequirement)}
}

// MultipartFileUploadTaskRequirement is the requirement for execute MultipartFileUploadTask.
type MultipartFileUploadTaskRequirement interface {
	navvy.Task
	types.Todoist
	types.PoolGetter

	// Inherited value
	types.KeyGetter
	types.MD5SumGetter
	types.OffsetGetter
	types.PartNumberGetter
	types.PathGetter
	types.SizeGetter
	types.StorageGetter
	types.UploadIDGetter
	types.WaitGroupGetter

	// Runtime value
}

// MultipartFileUploadTask will upload a multipart via file.
type MultipartFileUploadTask struct {
	MultipartFileUploadTaskRequirement
}

// mockMultipartFileUploadTask is the mock task for MultipartFileUploadTask.
type mockMultipartFileUploadTask struct {
	types.Todo
	types.Pool

	// Inherited value
	types.Key
	types.MD5Sum
	types.Offset
	types.PartNumber
	types.Path
	types.Size
	types.Storage
	types.UploadID
	types.WaitGroup

	// Runtime value
}

func (t *mockMultipartFileUploadTask) Run() {
	panic("mockMultipartFileUploadTask should not be run.")
}

// NewMultipartFileUploadTask will create a new MultipartFileUploadTask.
func NewMultipartFileUploadTask(task types.Todoist) navvy.Task {
	return &MultipartFileUploadTask{task.(MultipartFileUploadTaskRequirement)}
}

// MultipartInitTaskRequirement is the requirement for execute MultipartInitTask.
type MultipartInitTaskRequirement interface {
	navvy.Task
	types.Todoist
	types.PoolGetter

	// Inherited value
	types.CurrentOffsetGetter
	types.KeyGetter
	types.SizeGetter
	types.StorageGetter
	types.TaskConstructorGetter
	types.WaitGroupGetter

	// Runtime value
	types.UploadIDSetter
}

// MultipartInitTask will init a multipart upload.
type MultipartInitTask struct {
	MultipartInitTaskRequirement
}

// mockMultipartInitTask is the mock task for MultipartInitTask.
type mockMultipartInitTask struct {
	types.Todo
	types.Pool

	// Inherited value
	types.CurrentOffset
	types.Key
	types.Size
	types.Storage
	types.TaskConstructor
	types.WaitGroup

	// Runtime value
	types.UploadID
}

func (t *mockMultipartInitTask) Run() {
	panic("mockMultipartInitTask should not be run.")
}

// NewMultipartInitTask will create a new MultipartInitTask.
func NewMultipartInitTask(task types.Todoist) navvy.Task {
	return &MultipartInitTask{task.(MultipartInitTaskRequirement)}
}

// MultipartStreamUploadTaskRequirement is the requirement for execute MultipartStreamUploadTask.
type MultipartStreamUploadTaskRequirement interface {
	navvy.Task
	types.Todoist
	types.PoolGetter

	// Inherited value
	types.ContentGetter
	types.KeyGetter
	types.MD5SumGetter
	types.PartNumberGetter
	types.SizeGetter
	types.StorageGetter
	types.UploadIDGetter
	types.WaitGroupGetter

	// Runtime value
}

// MultipartStreamUploadTask will upload a multipart via stream.
type MultipartStreamUploadTask struct {
	MultipartStreamUploadTaskRequirement
}

// mockMultipartStreamUploadTask is the mock task for MultipartStreamUploadTask.
type mockMultipartStreamUploadTask struct {
	types.Todo
	types.Pool

	// Inherited value
	types.Content
	types.Key
	types.MD5Sum
	types.PartNumber
	types.Size
	types.Storage
	types.UploadID
	types.WaitGroup

	// Runtime value
}

func (t *mockMultipartStreamUploadTask) Run() {
	panic("mockMultipartStreamUploadTask should not be run.")
}

// NewMultipartStreamUploadTask will create a new MultipartStreamUploadTask.
func NewMultipartStreamUploadTask(task types.Todoist) navvy.Task {
	return &MultipartStreamUploadTask{task.(MultipartStreamUploadTaskRequirement)}
}

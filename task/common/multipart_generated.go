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

// multipartCompleteTaskRequirement is the requirement for execute MultipartCompleteTask.
type multipartCompleteTaskRequirement interface {
	navvy.Task
	types.Todoist
	types.PoolGetter
	types.FaultSetter
	types.FaultValidator
	types.IDGetter

	// Inherited value
	types.CurrentPartNumberGetter
	types.KeyGetter
	types.StorageGetter
	types.UploadIDGetter
	// Runtime value
}

// mockMultipartCompleteTask is the mock task for MultipartCompleteTask.
type mockMultipartCompleteTask struct {
	types.Todo
	types.Pool
	types.Fault
	types.ID

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

// MultipartCompleteTask will upload a multipart via stream.
type MultipartCompleteTask struct {
	multipartCompleteTaskRequirement
}

// Run implement navvy.Task.
func (t *MultipartCompleteTask) Run() {
	t.run()
	if t.ValidateFault() {
		return
	}
	utils.SubmitNextTask(t.multipartCompleteTaskRequirement)
}

func (t *MultipartCompleteTask) TriggerError(err error) {
	t.SetFault(fmt.Errorf("Task MultipartComplete failed: {%w}", err))
}

// NewMultipartCompleteTask will create a new MultipartCompleteTask.
func NewMultipartCompleteTask(task types.Todoist) navvy.Task {
	return &MultipartCompleteTask{task.(multipartCompleteTaskRequirement)}
}

// multipartFileUploadTaskRequirement is the requirement for execute MultipartFileUploadTask.
type multipartFileUploadTaskRequirement interface {
	navvy.Task
	types.Todoist
	types.PoolGetter
	types.FaultSetter
	types.FaultValidator
	types.IDGetter

	// Inherited value
	types.KeyGetter
	types.MD5SumGetter
	types.OffsetGetter
	types.PartNumberGetter
	types.PathGetter
	types.SchedulerGetter
	types.SizeGetter
	types.StorageGetter
	types.UploadIDGetter
	// Runtime value
}

// mockMultipartFileUploadTask is the mock task for MultipartFileUploadTask.
type mockMultipartFileUploadTask struct {
	types.Todo
	types.Pool
	types.Fault
	types.ID

	// Inherited value
	types.Key
	types.MD5Sum
	types.Offset
	types.PartNumber
	types.Path
	types.Scheduler
	types.Size
	types.Storage
	types.UploadID
	// Runtime value
}

func (t *mockMultipartFileUploadTask) Run() {
	panic("mockMultipartFileUploadTask should not be run.")
}

// MultipartFileUploadTask will upload a multipart via file.
type MultipartFileUploadTask struct {
	multipartFileUploadTaskRequirement
}

// Run implement navvy.Task.
func (t *MultipartFileUploadTask) Run() {
	t.run()
	if t.ValidateFault() {
		return
	}
	utils.SubmitNextTask(t.multipartFileUploadTaskRequirement)
}

func (t *MultipartFileUploadTask) TriggerError(err error) {
	t.SetFault(fmt.Errorf("Task MultipartFileUpload failed: {%w}", err))
}

// NewMultipartFileUploadTask will create a new MultipartFileUploadTask.
func NewMultipartFileUploadTask(task types.Todoist) navvy.Task {
	return &MultipartFileUploadTask{task.(multipartFileUploadTaskRequirement)}
}

// multipartInitTaskRequirement is the requirement for execute MultipartInitTask.
type multipartInitTaskRequirement interface {
	navvy.Task
	types.Todoist
	types.PoolGetter
	types.FaultSetter
	types.FaultValidator
	types.IDGetter

	// Inherited value
	types.CurrentOffsetGetter
	types.KeyGetter
	types.SchedulerGetter
	types.StorageGetter
	types.TotalSizeGetter
	// Runtime value
	types.UploadIDSetter
}

// mockMultipartInitTask is the mock task for MultipartInitTask.
type mockMultipartInitTask struct {
	types.Todo
	types.Pool
	types.Fault
	types.ID

	// Inherited value
	types.CurrentOffset
	types.Key
	types.Scheduler
	types.Storage
	types.TotalSize
	// Runtime value
	types.UploadID
}

func (t *mockMultipartInitTask) Run() {
	panic("mockMultipartInitTask should not be run.")
}

// MultipartInitTask will init a multipart upload.
type MultipartInitTask struct {
	multipartInitTaskRequirement
}

// Run implement navvy.Task.
func (t *MultipartInitTask) Run() {
	t.run()
	if t.ValidateFault() {
		return
	}
	utils.SubmitNextTask(t.multipartInitTaskRequirement)
}

func (t *MultipartInitTask) TriggerError(err error) {
	t.SetFault(fmt.Errorf("Task MultipartInit failed: {%w}", err))
}

// NewMultipartInitTask will create a new MultipartInitTask.
func NewMultipartInitTask(task types.Todoist) navvy.Task {
	return &MultipartInitTask{task.(multipartInitTaskRequirement)}
}

// multipartStreamUploadTaskRequirement is the requirement for execute MultipartStreamUploadTask.
type multipartStreamUploadTaskRequirement interface {
	navvy.Task
	types.Todoist
	types.PoolGetter
	types.FaultSetter
	types.FaultValidator
	types.IDGetter

	// Inherited value
	types.ContentGetter
	types.KeyGetter
	types.MD5SumGetter
	types.PartNumberGetter
	types.SchedulerGetter
	types.SizeGetter
	types.StorageGetter
	types.UploadIDGetter
	// Runtime value
}

// mockMultipartStreamUploadTask is the mock task for MultipartStreamUploadTask.
type mockMultipartStreamUploadTask struct {
	types.Todo
	types.Pool
	types.Fault
	types.ID

	// Inherited value
	types.Content
	types.Key
	types.MD5Sum
	types.PartNumber
	types.Scheduler
	types.Size
	types.Storage
	types.UploadID
	// Runtime value
}

func (t *mockMultipartStreamUploadTask) Run() {
	panic("mockMultipartStreamUploadTask should not be run.")
}

// MultipartStreamUploadTask will upload a multipart via stream.
type MultipartStreamUploadTask struct {
	multipartStreamUploadTaskRequirement
}

// Run implement navvy.Task.
func (t *MultipartStreamUploadTask) Run() {
	t.run()
	if t.ValidateFault() {
		return
	}
	utils.SubmitNextTask(t.multipartStreamUploadTaskRequirement)
}

func (t *MultipartStreamUploadTask) TriggerError(err error) {
	t.SetFault(fmt.Errorf("Task MultipartStreamUpload failed: {%w}", err))
}

// NewMultipartStreamUploadTask will create a new MultipartStreamUploadTask.
func NewMultipartStreamUploadTask(task types.Todoist) navvy.Task {
	return &MultipartStreamUploadTask{task.(multipartStreamUploadTaskRequirement)}
}

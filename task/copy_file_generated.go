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

// copyFileTaskRequirement is the requirement for execute CopyFileTask.
type copyFileTaskRequirement interface {
	navvy.Task
	types.PoolGetter

	// Inherited value
	types.KeyGetter
	types.PathGetter
	types.StorageGetter
}

// mockCopyFileTask is the mock task for CopyFileTask.
type mockCopyFileTask struct {
	types.Todo
	types.Pool
	types.Fault
	types.ID

	// Inherited value
	types.Key
	types.Path
	types.Storage
}

func (t *mockCopyFileTask) Run() {
	panic("mockCopyFileTask should not be run.")
}

// CopyFileTask will copy a file to storage.
type CopyFileTask struct {
	copyFileTaskRequirement

	// Predefined runtime value
	types.Fault
	types.ID
	types.Todo

	// Runtime value
	types.TotalSize
}

// Run implement navvy.Task
func (t *CopyFileTask) Run() {
	if t.ValidateFault() {
		return
	}
	utils.SubmitNextTask(t)
}

func (t *CopyFileTask) TriggerError(err error) {
	t.SetFault(fmt.Errorf("Task CopyFile failed: {%w}", err))
}

// NewCopyFileTask will create a CopyFileTask and fetch inherited data from CopyTask.
func NewCopyFileTask(task types.Todoist) navvy.Task {
	t := &CopyFileTask{
		copyFileTaskRequirement: task.(copyFileTaskRequirement),
	}
	t.SetID(uuid.New().String())
	t.new()
	return t
}

// copyLargeFileTaskRequirement is the requirement for execute CopyLargeFileTask.
type copyLargeFileTaskRequirement interface {
	navvy.Task
	types.PoolGetter

	// Inherited value
	types.KeyGetter
	types.PathGetter
	types.StorageGetter
	types.TotalSizeGetter
}

// mockCopyLargeFileTask is the mock task for CopyLargeFileTask.
type mockCopyLargeFileTask struct {
	types.Todo
	types.Pool
	types.Fault
	types.ID

	// Inherited value
	types.Key
	types.Path
	types.Storage
	types.TotalSize
}

func (t *mockCopyLargeFileTask) Run() {
	panic("mockCopyLargeFileTask should not be run.")
}

// CopyLargeFileTask will copy a large file to storage.
type CopyLargeFileTask struct {
	copyLargeFileTaskRequirement

	// Predefined runtime value
	types.Fault
	types.ID
	types.Todo

	// Runtime value
	types.CurrentOffset
	types.CurrentPartNumber
	types.PartSize
	types.Scheduler
	types.UploadID
}

// Run implement navvy.Task
func (t *CopyLargeFileTask) Run() {
	if t.ValidateFault() {
		return
	}
	utils.SubmitNextTask(t)
}

func (t *CopyLargeFileTask) TriggerError(err error) {
	t.SetFault(fmt.Errorf("Task CopyLargeFile failed: {%w}", err))
}

// NewCopyLargeFileTask will create a CopyLargeFileTask and fetch inherited data from CopyFileTask.
func NewCopyLargeFileTask(task types.Todoist) navvy.Task {
	t := &CopyLargeFileTask{
		copyLargeFileTaskRequirement: task.(copyLargeFileTaskRequirement),
	}
	t.SetID(uuid.New().String())
	t.new()
	return t
}

// copyPartialFileTaskRequirement is the requirement for execute CopyPartialFileTask.
type copyPartialFileTaskRequirement interface {
	navvy.Task
	types.PoolGetter

	// Inherited value
	types.CurrentOffsetGetter
	types.CurrentPartNumberGetter
	types.KeyGetter
	types.PartSizeGetter
	types.PathGetter
	types.SchedulerGetter
	types.StorageGetter
	types.TotalSizeGetter
	types.UploadIDGetter
}

// mockCopyPartialFileTask is the mock task for CopyPartialFileTask.
type mockCopyPartialFileTask struct {
	types.Todo
	types.Pool
	types.Fault
	types.ID

	// Inherited value
	types.CurrentOffset
	types.CurrentPartNumber
	types.Key
	types.PartSize
	types.Path
	types.Scheduler
	types.Storage
	types.TotalSize
	types.UploadID
}

func (t *mockCopyPartialFileTask) Run() {
	panic("mockCopyPartialFileTask should not be run.")
}

// CopyPartialFileTask will copy a partial file to storage, is the sub task for CopyLargeFile.
type CopyPartialFileTask struct {
	copyPartialFileTaskRequirement

	// Predefined runtime value
	types.Fault
	types.ID
	types.Todo

	// Runtime value
	types.MD5Sum
	types.Offset
	types.PartNumber
	types.Size
}

// Run implement navvy.Task
func (t *CopyPartialFileTask) Run() {
	if t.ValidateFault() {
		return
	}
	utils.SubmitNextTask(t)
}

func (t *CopyPartialFileTask) TriggerError(err error) {
	t.SetFault(fmt.Errorf("Task CopyPartialFile failed: {%w}", err))
}

// NewCopyPartialFileTask will create a CopyPartialFileTask and fetch inherited data from CopyLargeFileTask.
func NewCopyPartialFileTask(task types.Todoist) navvy.Task {
	t := &CopyPartialFileTask{
		copyPartialFileTaskRequirement: task.(copyPartialFileTaskRequirement),
	}
	t.SetID(uuid.New().String())
	t.new()
	return t
}

// copySmallFileTaskRequirement is the requirement for execute CopySmallFileTask.
type copySmallFileTaskRequirement interface {
	navvy.Task
	types.PoolGetter

	// Inherited value
	types.KeyGetter
	types.PathGetter
	types.StorageGetter
	types.TotalSizeGetter
}

// mockCopySmallFileTask is the mock task for CopySmallFileTask.
type mockCopySmallFileTask struct {
	types.Todo
	types.Pool
	types.Fault
	types.ID

	// Inherited value
	types.Key
	types.Path
	types.Storage
	types.TotalSize
}

func (t *mockCopySmallFileTask) Run() {
	panic("mockCopySmallFileTask should not be run.")
}

// CopySmallFileTask will copy a small file to storage.
type CopySmallFileTask struct {
	copySmallFileTaskRequirement

	// Predefined runtime value
	types.Fault
	types.ID
	types.Todo

	// Runtime value
	types.MD5Sum
	types.Offset
	types.Size
}

// Run implement navvy.Task
func (t *CopySmallFileTask) Run() {
	if t.ValidateFault() {
		return
	}
	utils.SubmitNextTask(t)
}

func (t *CopySmallFileTask) TriggerError(err error) {
	t.SetFault(fmt.Errorf("Task CopySmallFile failed: {%w}", err))
}

// NewCopySmallFileTask will create a CopySmallFileTask and fetch inherited data from CopyFileTask.
func NewCopySmallFileTask(task types.Todoist) navvy.Task {
	t := &CopySmallFileTask{
		copySmallFileTaskRequirement: task.(copySmallFileTaskRequirement),
	}
	t.SetID(uuid.New().String())
	t.new()
	return t
}

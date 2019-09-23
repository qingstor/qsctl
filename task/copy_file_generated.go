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

// copyFileTaskRequirement is the requirement for execute CopyFileTask.
type copyFileTaskRequirement interface {
	navvy.Task
	types.PoolGetter

	// Inherited value
	types.KeyGetter
	types.PathGetter
	types.StorageGetter
}

// CopyFileTask will copy a file to storage.
type CopyFileTask struct {
	copyFileTaskRequirement

	// Runtime value
	types.Todo
	types.TotalSize
}

// mockCopyFileTask is the mock task for CopyFileTask.
type mockCopyFileTask struct {
	types.Todo
	types.Pool

	// Inherited value
	types.Key
	types.Path
	types.Storage
}

func (t *mockCopyFileTask) Run() {
	panic("mockCopyFileTask should not be run.")
}

// Run implement navvy.Task
func (t *CopyFileTask) Run() {
	utils.SubmitNextTask(t)
}

// initCopyFileTask will create a CopyFileTask and fetch inherited data from CopyTask.
func NewCopyFileTask(task types.Todoist) navvy.Task {
	t := &CopyFileTask{
		copyFileTaskRequirement: task.(copyFileTaskRequirement),
	}
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

// CopyLargeFileTask will copy a large file to storage.
type CopyLargeFileTask struct {
	copyLargeFileTaskRequirement

	// Runtime value
	types.Todo
	types.CurrentOffset
	types.CurrentPartNumber
	types.PartSize
	types.TaskConstructor
	types.UploadID
	types.WaitGroup
}

// mockCopyLargeFileTask is the mock task for CopyLargeFileTask.
type mockCopyLargeFileTask struct {
	types.Todo
	types.Pool

	// Inherited value
	types.Key
	types.Path
	types.Storage
	types.TotalSize
}

func (t *mockCopyLargeFileTask) Run() {
	panic("mockCopyLargeFileTask should not be run.")
}

// Run implement navvy.Task
func (t *CopyLargeFileTask) Run() {
	utils.SubmitNextTask(t)
}

// initCopyLargeFileTask will create a CopyLargeFileTask and fetch inherited data from CopyFileTask.
func NewCopyLargeFileTask(task types.Todoist) navvy.Task {
	t := &CopyLargeFileTask{
		copyLargeFileTaskRequirement: task.(copyLargeFileTaskRequirement),
	}
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
	types.StorageGetter
	types.TotalSizeGetter
	types.UploadIDGetter
	types.WaitGroupGetter
}

// CopyPartialFileTask will copy a partial file to storage, is the sub task for CopyLargeFile.
type CopyPartialFileTask struct {
	copyPartialFileTaskRequirement

	// Runtime value
	types.Todo
	types.MD5Sum
	types.Offset
	types.PartNumber
	types.Size
}

// mockCopyPartialFileTask is the mock task for CopyPartialFileTask.
type mockCopyPartialFileTask struct {
	types.Todo
	types.Pool

	// Inherited value
	types.CurrentOffset
	types.CurrentPartNumber
	types.Key
	types.PartSize
	types.Path
	types.Storage
	types.TotalSize
	types.UploadID
	types.WaitGroup
}

func (t *mockCopyPartialFileTask) Run() {
	panic("mockCopyPartialFileTask should not be run.")
}

// Run implement navvy.Task
func (t *CopyPartialFileTask) Run() {
	utils.SubmitNextTask(t)
}

// initCopyPartialFileTask will create a CopyPartialFileTask and fetch inherited data from CopyLargeFileTask.
func NewCopyPartialFileTask(task types.Todoist) navvy.Task {
	t := &CopyPartialFileTask{
		copyPartialFileTaskRequirement: task.(copyPartialFileTaskRequirement),
	}
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

// CopySmallFileTask will copy a small file to storage.
type CopySmallFileTask struct {
	copySmallFileTaskRequirement

	// Runtime value
	types.Todo
	types.MD5Sum
	types.Offset
	types.Size
}

// mockCopySmallFileTask is the mock task for CopySmallFileTask.
type mockCopySmallFileTask struct {
	types.Todo
	types.Pool

	// Inherited value
	types.Key
	types.Path
	types.Storage
	types.TotalSize
}

func (t *mockCopySmallFileTask) Run() {
	panic("mockCopySmallFileTask should not be run.")
}

// Run implement navvy.Task
func (t *CopySmallFileTask) Run() {
	utils.SubmitNextTask(t)
}

// initCopySmallFileTask will create a CopySmallFileTask and fetch inherited data from CopyFileTask.
func NewCopySmallFileTask(task types.Todoist) navvy.Task {
	t := &CopySmallFileTask{
		copySmallFileTaskRequirement: task.(copySmallFileTaskRequirement),
	}
	t.new()
	return t
}

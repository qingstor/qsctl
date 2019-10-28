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

// fileMD5SumTaskRequirement is the requirement for execute FileMD5SumTask.
type fileMD5SumTaskRequirement interface {
	navvy.Task

	// Inherited value
	types.OffsetGetter
	types.SizeGetter
	types.SourcePathGetter
	types.SourceStorageGetter

	// Mutable value
	types.MD5SumSetter
}

// mockFileMD5SumTask is the mock task for FileMD5SumTask.
type mockFileMD5SumTask struct {
	types.Pool
	types.Fault
	types.ID

	// Inherited value
	types.Offset
	types.Size
	types.SourcePath
	types.SourceStorage

	// Mutable value
	types.MD5Sum
}

func (t *mockFileMD5SumTask) Run() {
	panic("mockFileMD5SumTask should not be run.")
}

// FileMD5SumTask will get file's md5 sum.
type FileMD5SumTask struct {
	fileMD5SumTaskRequirement

	// Predefined runtime value
	types.Fault
	types.ID
	types.Scheduler

	// Runtime value
}

// Run implement navvy.Task
func (t *FileMD5SumTask) Run() {
	t.run()
}

func (t *FileMD5SumTask) TriggerFault(err error) {
	t.SetFault(fmt.Errorf("Task FileMD5Sum failed: {%w}", err))
}

// NewFileMD5SumTask will create a FileMD5SumTask and fetch inherited data from parent task.
func NewFileMD5SumTask(task navvy.Task) navvy.Task {
	t := &FileMD5SumTask{
		fileMD5SumTaskRequirement: task.(fileMD5SumTaskRequirement),
	}
	t.SetID(uuid.New().String())
	t.new()
	return t
}

// streamMD5SumTaskRequirement is the requirement for execute StreamMD5SumTask.
type streamMD5SumTaskRequirement interface {
	navvy.Task

	// Inherited value
	types.ContentGetter

	// Mutable value
}

// mockStreamMD5SumTask is the mock task for StreamMD5SumTask.
type mockStreamMD5SumTask struct {
	types.Pool
	types.Fault
	types.ID

	// Inherited value
	types.Content

	// Mutable value
}

func (t *mockStreamMD5SumTask) Run() {
	panic("mockStreamMD5SumTask should not be run.")
}

// StreamMD5SumTask will get stream's md5 sum.
type StreamMD5SumTask struct {
	streamMD5SumTaskRequirement

	// Predefined runtime value
	types.Fault
	types.ID
	types.Scheduler

	// Runtime value
	types.MD5Sum
}

// Run implement navvy.Task
func (t *StreamMD5SumTask) Run() {
	t.run()
}

func (t *StreamMD5SumTask) TriggerFault(err error) {
	t.SetFault(fmt.Errorf("Task StreamMD5Sum failed: {%w}", err))
}

// NewStreamMD5SumTask will create a StreamMD5SumTask and fetch inherited data from parent task.
func NewStreamMD5SumTask(task navvy.Task) navvy.Task {
	t := &StreamMD5SumTask{
		streamMD5SumTaskRequirement: task.(streamMD5SumTaskRequirement),
	}
	t.SetID(uuid.New().String())
	t.new()
	return t
}

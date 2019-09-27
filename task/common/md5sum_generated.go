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

// fileMD5SumTaskRequirement is the requirement for execute FileMD5SumTask.
type fileMD5SumTaskRequirement interface {
	navvy.Task
	types.Todoist
	types.PoolGetter
	types.FaultSetter
	types.FaultValidator
	types.IDGetter

	// Inherited value
	types.OffsetGetter
	types.PathGetter
	types.SizeGetter
	// Runtime value
	types.MD5SumSetter
}

// mockFileMD5SumTask is the mock task for FileMD5SumTask.
type mockFileMD5SumTask struct {
	types.Todo
	types.Pool
	types.Fault
	types.ID

	// Inherited value
	types.Offset
	types.Path
	types.Size
	// Runtime value
	types.MD5Sum
}

func (t *mockFileMD5SumTask) Run() {
	panic("mockFileMD5SumTask should not be run.")
}

// FileMD5SumTask will get file's md5 sum.
type FileMD5SumTask struct {
	fileMD5SumTaskRequirement
}

// Run implement navvy.Task.
func (t *FileMD5SumTask) Run() {
	t.run()
	if t.ValidateFault() {
		return
	}
	utils.SubmitNextTask(t.fileMD5SumTaskRequirement)
}

func (t *FileMD5SumTask) TriggerError(err error) {
	t.SetFault(fmt.Errorf("Task FileMD5Sum failed: {%w}", err))
}

// NewFileMD5SumTask will create a new FileMD5SumTask.
func NewFileMD5SumTask(task types.Todoist) navvy.Task {
	return &FileMD5SumTask{task.(fileMD5SumTaskRequirement)}
}

// streamMD5SumTaskRequirement is the requirement for execute StreamMD5SumTask.
type streamMD5SumTaskRequirement interface {
	navvy.Task
	types.Todoist
	types.PoolGetter
	types.FaultSetter
	types.FaultValidator
	types.IDGetter

	// Inherited value
	types.ContentGetter
	// Runtime value
	types.MD5SumSetter
}

// mockStreamMD5SumTask is the mock task for StreamMD5SumTask.
type mockStreamMD5SumTask struct {
	types.Todo
	types.Pool
	types.Fault
	types.ID

	// Inherited value
	types.Content
	// Runtime value
	types.MD5Sum
}

func (t *mockStreamMD5SumTask) Run() {
	panic("mockStreamMD5SumTask should not be run.")
}

// StreamMD5SumTask will get stream's md5 sum.
type StreamMD5SumTask struct {
	streamMD5SumTaskRequirement
}

// Run implement navvy.Task.
func (t *StreamMD5SumTask) Run() {
	t.run()
	if t.ValidateFault() {
		return
	}
	utils.SubmitNextTask(t.streamMD5SumTaskRequirement)
}

func (t *StreamMD5SumTask) TriggerError(err error) {
	t.SetFault(fmt.Errorf("Task StreamMD5Sum failed: {%w}", err))
}

// NewStreamMD5SumTask will create a new StreamMD5SumTask.
func NewStreamMD5SumTask(task types.Todoist) navvy.Task {
	return &StreamMD5SumTask{task.(streamMD5SumTaskRequirement)}
}

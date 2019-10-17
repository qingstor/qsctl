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

// copyPartialStreamTaskRequirement is the requirement for execute CopyPartialStreamTask.
type copyPartialStreamTaskRequirement interface {
	navvy.Task
	types.PoolGetter

	// Inherited value
	types.BytesPoolGetter
	types.CurrentOffsetGetter
	types.DestinationStorageGetter
	types.KeyGetter
	types.PartSizeGetter
	types.SchedulerGetter
	types.SegmentIDGetter
	types.StreamGetter
}

// mockCopyPartialStreamTask is the mock task for CopyPartialStreamTask.
type mockCopyPartialStreamTask struct {
	types.Todo
	types.Pool
	types.Fault
	types.ID

	// Inherited value
	types.BytesPool
	types.CurrentOffset
	types.DestinationStorage
	types.Key
	types.PartSize
	types.Scheduler
	types.SegmentID
	types.Stream
}

func (t *mockCopyPartialStreamTask) Run() {
	panic("mockCopyPartialStreamTask should not be run.")
}

// CopyPartialStreamTask will copy a partial stream to DestinationStorage.
type CopyPartialStreamTask struct {
	copyPartialStreamTaskRequirement

	// Predefined runtime value
	types.Fault
	types.ID
	types.Todo

	// Runtime value
	types.Content
	types.MD5Sum
	types.Offset
	types.Size
}

// Run implement navvy.Task
func (t *CopyPartialStreamTask) Run() {
	if t.ValidateFault() {
		return
	}
	utils.SubmitNextTask(t)
}

func (t *CopyPartialStreamTask) TriggerFault(err error) {
	t.SetFault(fmt.Errorf("Task CopyPartialStream failed: {%w}", err))
}
// NewCopyPartialStreamTask will create a CopyPartialStreamTask and fetch inherited data from CopyStreamTask.
func NewCopyPartialStreamTask(task types.Todoist) navvy.Task {
	t := &CopyPartialStreamTask{
		copyPartialStreamTaskRequirement: task.(copyPartialStreamTaskRequirement),
	}
	t.SetID(uuid.New().String())
	t.new()
	return t
}

// copyStreamTaskRequirement is the requirement for execute CopyStreamTask.
type copyStreamTaskRequirement interface {
	navvy.Task
	types.PoolGetter

	// Inherited value
	types.DestinationStorageGetter
	types.KeyGetter
	types.StreamGetter
}

// mockCopyStreamTask is the mock task for CopyStreamTask.
type mockCopyStreamTask struct {
	types.Todo
	types.Pool
	types.Fault
	types.ID

	// Inherited value
	types.DestinationStorage
	types.Key
	types.Stream
}

func (t *mockCopyStreamTask) Run() {
	panic("mockCopyStreamTask should not be run.")
}

// CopyStreamTask will copy a stream to DestinationStorage.
type CopyStreamTask struct {
	copyStreamTaskRequirement

	// Predefined runtime value
	types.Fault
	types.ID
	types.Todo

	// Runtime value
	types.BytesPool
	types.CurrentOffset
	types.PartSize
	types.Scheduler
	types.SegmentID
	types.TotalSize
}

// Run implement navvy.Task
func (t *CopyStreamTask) Run() {
	if t.ValidateFault() {
		return
	}
	utils.SubmitNextTask(t)
}

func (t *CopyStreamTask) TriggerFault(err error) {
	t.SetFault(fmt.Errorf("Task CopyStream failed: {%w}", err))
}
// NewCopyStreamTask will create a CopyStreamTask and fetch inherited data from CopyTask.
func NewCopyStreamTask(task types.Todoist) navvy.Task {
	t := &CopyStreamTask{
		copyStreamTaskRequirement: task.(copyStreamTaskRequirement),
	}
	t.SetID(uuid.New().String())
	t.new()
	return t
}

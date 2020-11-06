package taskutils

import (
	"github.com/qingstor/noah/pkg/task"

	"github.com/qingstor/noah/pkg/types"
)

// AtServiceTask is the root task for service only task.
type AtServiceTask struct {
	task.Task

	types.Service
}

// NewAtServiceTask will create a new between storage task.
func NewAtServiceTask() *AtServiceTask {
	t := &AtServiceTask{}
	return t
}

// AtStorageTask is the root task for single storage task.
type AtStorageTask struct {
	task.Task

	types.Path
	types.Storage
	types.Type
}

// NewAtStorageTask will create a new between storage task.
func NewAtStorageTask() *AtStorageTask {
	t := &AtStorageTask{}
	return t
}

// BetweenStorageTask is the root task for tasks operate between two storager.
type BetweenStorageTask struct {
	task.Task

	types.SourcePath
	types.SourceStorage
	types.SourceType
	types.DestinationPath
	types.DestinationStorage
	types.DestinationType
}

// NewBetweenStorageTask will create a new between storage task.
func NewBetweenStorageTask() *BetweenStorageTask {
	t := &BetweenStorageTask{}
	return t
}

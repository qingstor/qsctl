package taskutils

import (
	"github.com/Xuanwo/navvy"

	"github.com/qingstor/noah/pkg/fault"
	"github.com/qingstor/noah/pkg/types"
)

// AtServiceTask is the root task for service only task.
type AtServiceTask struct {
	navvy.Task
	types.Pool
	types.Fault

	types.Service
}

// NewAtServiceTask will create a new between storage task.
func NewAtServiceTask(poolSize int) *AtServiceTask {
	t := &AtServiceTask{}
	t.SetPool(navvy.NewPool(poolSize))
	t.SetFault(fault.New())
	return t
}

// AtStorageTask is the root task for single storage task.
type AtStorageTask struct {
	navvy.Task
	types.Pool
	types.Fault

	types.WorkDir
	types.Path
	types.Storage
	types.Type
}

// NewAtStorageTask will create a new between storage task.
func NewAtStorageTask(poolSize int) *AtStorageTask {
	t := &AtStorageTask{}
	t.SetPool(navvy.NewPool(poolSize))
	t.SetFault(fault.New())
	return t
}

// BetweenStorageTask is the root task for tasks operate between two storager.
type BetweenStorageTask struct {
	navvy.Task
	types.Pool
	types.Fault

	types.SourceWorkDir
	types.SourcePath
	types.SourceStorage
	types.SourceType
	types.DestinationWorkDir
	types.DestinationPath
	types.DestinationStorage
	types.DestinationType
}

// NewBetweenStorageTask will create a new between storage task.
func NewBetweenStorageTask(poolSize int) *BetweenStorageTask {
	t := &BetweenStorageTask{}
	t.SetPool(navvy.NewPool(poolSize))
	t.SetFault(fault.New())
	return t
}

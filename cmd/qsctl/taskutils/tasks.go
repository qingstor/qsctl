package taskutils

import (
	"github.com/Xuanwo/navvy"
	"github.com/yunify/qsctl/v2/pkg/types"
)

// AtServiceTask is the root task for service only task.
type AtServiceTask struct {
	types.Service
}

// AtStorageTask is the root task for single storage task.
type AtStorageTask struct {
	types.Path
	types.Storage
	types.Type
}

// BetweenStorageTask is the root task for tasks operate between two storager.
type BetweenStorageTask struct {
	navvy.Task
	types.Pool

	types.SourcePath
	types.SourceStorage
	types.SourceType
	types.DestinationPath
	types.DestinationStorage
	types.DestinationType
}

func NewBetweenStorageTask(poolSize int) *BetweenStorageTask {
	t := &BetweenStorageTask{}
	t.SetPool(navvy.NewPool(poolSize))
	return t
}

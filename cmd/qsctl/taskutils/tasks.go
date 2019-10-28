package taskutils

import (
	"github.com/yunify/qsctl/v2/pkg/types"
)

type AtServiceTask struct {
	types.Service
}

type AtStorageTask struct {
	types.Path
	types.Storage
	types.Type
}

type BetweenStorageTask struct {
	types.SourcePath
	types.SourceStorage
	types.SourceType
	types.DestinationPath
	types.DestinationStorage
	types.DestinationType
}

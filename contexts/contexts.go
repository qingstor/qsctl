package contexts

import (
	"github.com/Xuanwo/navvy"
	"github.com/yunify/qsctl/v2/storage"
)

var (
	// Storage is the remote storage service.
	// TODO: move to action local var.
	Storage storage.ObjectStorage
	// Pool is the task pool.
	Pool *navvy.Pool
)

// SetupServices will setup services.
func SetupServices() (err error) {
	Storage, err = storage.NewQingStorObjectStorage()
	if err != nil {
		return
	}

	Pool, err = navvy.NewPool(10)
	if err != nil {
		return
	}
	return
}

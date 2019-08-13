package contexts

import (
	"github.com/yunify/qsctl/v2/storage"
)

var (
	// Storage is the remote storage service.
	Storage storage.ObjectStorage
)

// SetupServices will setup services.
func SetupServices() (err error) {
	Storage, err = storage.NewQingStorObjectStorage()
	if err != nil {
		return
	}

	return
}

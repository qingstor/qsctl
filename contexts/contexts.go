package contexts

import (
	"github.com/yunify/qsctl/storage"
)

var (
	// Storage is the remote storage service.
	Storage storage.ObjectStorage
)

// Available flags.
var (
	// Global flags.
	Bench bool
	// Copy commands flags.
	ExpectSize           int64
	MaximumMemoryContent int64
	// Bucket operation flags.
	Zone string
)

// SetupServices will setup services.
func SetupServices() (err error) {
	Storage, err = storage.NewQingStorObjectStorage()
	if err != nil {
		return
	}

	return
}

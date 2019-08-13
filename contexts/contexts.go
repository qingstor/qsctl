package contexts

import (
	"context"

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

// SetContext set k-v into specific ctx
func SetContext(ctx context.Context, k, v interface{}) context.Context {
	return context.WithValue(ctx, k, v)
}

// FromContext get value from context with the specific k
func FromContext(ctx context.Context, k interface{}) interface{} {
	return ctx.Value(k)
}

// NewCmdCtx returns a new CmdCtx with empty background
func NewCmdCtx() context.Context {
	return context.Background()
}

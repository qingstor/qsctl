package utils

import (
	"github.com/c2h5oh/datasize"
	"github.com/google/uuid"
	"github.com/yunify/qsctl/v2/pkg/types"
)

const maxTestSize = 64 * int64(datasize.MB)

var _ = uuid.New()

// EmptyTask is used for test.
type EmptyTask struct {
	types.ID
	types.Fault
	types.Pool
}

// Run implement navvy.Task interface.
func (t *EmptyTask) Run() {
}

// NewCallbackTask will create a new callback test.
func NewCallbackTask(fn func()) *CallbackTask {
	t := &CallbackTask{
		fn: fn,
	}
	t.SetID(uuid.New().String())
	return t
}

// CallbackTask is the callback task.
type CallbackTask struct {
	types.ID
	types.Fault
	types.Pool

	fn func()
}

// Run implement navvy.Task interface.
func (t *CallbackTask) Run() {
	t.fn()
}

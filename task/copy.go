package task

import (
	"github.com/Xuanwo/navvy"

	stypes "github.com/Xuanwo/storage/types"

	"github.com/yunify/qsctl/v2/pkg/types"
)

var copyTaskConstructor = map[stypes.ObjectType]types.TodoFunc{
	stypes.ObjectTypeStream: NewCopyStreamTask,
	stypes.ObjectTypeFile:   NewCopyFileTask,
}

// NewCopyTask will create a copy task.
func NewCopyTask(fn func(*CopyTask)) *CopyTask {
	t := &CopyTask{}

	pool := navvy.NewPool(10)
	t.SetPool(pool)

	fn(t)
	if t.ValidateFault() {
		return t
	}

	todo := copyTaskConstructor[t.GetSourceType()]
	if todo == nil {
		panic("invalid todo func")
	}
	t.AddTODOs(todo)
	return t
}

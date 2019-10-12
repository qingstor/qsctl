package task

import (
	"github.com/Xuanwo/navvy"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/pkg/types"
)

var copyTaskConstructor = map[constants.FlowType]map[constants.PathType]types.TodoFunc{
	constants.FlowToRemote: {
		constants.PathTypeStream: NewCopyStreamTask,
		constants.PathTypeFile:   NewCopyFileTask,
	},
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

	todo := copyTaskConstructor[t.GetFlowType()][t.GetPathType()]
	if todo == nil {
		panic("invalid todo func")
	}
	t.AddTODOs(todo)
	return t
}

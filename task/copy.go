package task

import (
	"github.com/Xuanwo/navvy"
	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/task/types"
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

	pool, err := navvy.NewPool(10)
	if err != nil {
		panic(err)
	}
	t.SetPool(pool)

	fn(t)

	todo := copyTaskConstructor[t.GetFlowType()][t.GetPathType()]
	if todo == nil {
		panic("invalid todo func")
	}
	t.AddTODOs(todo)
	return t
}

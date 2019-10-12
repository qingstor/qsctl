package task

import (
	"github.com/Xuanwo/navvy"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/pkg/types"
	"github.com/yunify/qsctl/v2/task/common"
)

var listTaskConstructor = map[constants.ListType]types.TodoFunc{
	constants.ListTypeBucket: common.NewBucketListTask,
}

// NewListTask will create a list task.
func NewListTask(fn func(*ListTask)) *ListTask {
	t := &ListTask{}

	pool := navvy.NewPool(10)
	t.SetPool(pool)

	fn(t)

	todo := listTaskConstructor[t.GetListType()]
	if todo == nil {
		panic("invalid todo func")
	}
	t.AddTODOs(todo)
	return t
}

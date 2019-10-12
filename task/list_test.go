package task

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/pkg/types"
	"github.com/yunify/qsctl/v2/task/common"
)

func TestNewListTask(t *testing.T) {
	cases := []struct {
		listType         constants.ListType
		expectedTodoFunc types.TodoFunc
		wantPanic        bool
	}{
		{constants.ListTypeBucket, common.NewBucketListTask, false},
		{constants.ListTypeInvalid, nil, true},
	}

	for _, v := range cases {
		pt := new(ListTask)
		panicFunc := func() {
			pt = NewListTask(func(task *ListTask) {
				task.SetListType(v.listType)
			})
		}
		if v.wantPanic {
			assert.Panics(t, panicFunc)
		} else {
			panicFunc()
		}
		assert.Equal(t,
			fmt.Sprintf("%v", v.expectedTodoFunc),
			fmt.Sprintf("%v", pt.Todo.NextTODO()))
	}
}

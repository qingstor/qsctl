package task

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yunify/qsctl/v2/pkg/types"
	"github.com/yunify/qsctl/v2/task/common"
)

func TestNewRemoveBucketTask(t *testing.T) {
	cases := []struct {
		input            string
		force            bool
		expectedTodoFunc types.TodoFunc
		expectErr        error
	}{
		{"qs://test-bucket/obj", false, common.NewBucketDeleteTask, nil},
		{"test-bucket", false, common.NewBucketDeleteTask, nil},
	}

	for _, v := range cases {
		pt := NewRemoveBucketTask(func(task *RemoveBucketTask) {
		})

		assert.Equal(t,
			fmt.Sprintf("%v", v.expectedTodoFunc),
			fmt.Sprintf("%v", pt.Todo.NextTODO()))
	}
}

package task

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yunify/qsctl/v2/pkg/types"
	"github.com/yunify/qsctl/v2/task/common"
	"github.com/yunify/qsctl/v2/utils"
)

func TestNewMakeBucketTask(t *testing.T) {
	cases := []struct {
		input            string
		expectedTodoFunc types.TodoFunc
		expectErr        error
	}{
		{"qs://test-bucket", common.NewBucketCreateTask, nil},
		{"test-bucket", common.NewBucketCreateTask, nil},
	}

	for _, v := range cases {
		pt := NewMakeBucketTask(func(task *MakeBucketTask) {
			_, bucketName, _, err := utils.ParseQsPath(v.input)
			if err != nil {
				t.Fatal(err)
			}
			task.SetBucketName(bucketName)
		})

		assert.Equal(t,
			fmt.Sprintf("%v", v.expectedTodoFunc),
			fmt.Sprintf("%v", pt.Todo.NextTODO()))
	}
}

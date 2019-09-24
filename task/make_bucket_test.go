package task

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/task/common"
	"github.com/yunify/qsctl/v2/task/types"
	"github.com/yunify/qsctl/v2/task/utils"
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
			keyType, bucketName, _, err := utils.ParseKey(v.input)
			if err != nil {
				t.Fatal(err)
			}
			if keyType != constants.KeyTypeBucket {
				t.Logf("key type: %d", keyType)
				t.Fatal(constants.ErrorQsPathInvalid)
			}
			task.SetBucketName(bucketName)
		})

		assert.Equal(t,
			fmt.Sprintf("%v", v.expectedTodoFunc),
			fmt.Sprintf("%v", pt.Todo.NextTODO()))
	}
}

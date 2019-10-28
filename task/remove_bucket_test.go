package task

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yunify/qsctl/v2/pkg/types"
	"github.com/yunify/qsctl/v2/task/common"
)

func TestNewRemoveBucketTask(t *testing.T) {
	removeBucketErr := errors.New("remove bucket error")
	cases := []struct {
		input            string
		force            bool
		expectedTodoFunc types.TaskFunc
		expectErr        error
	}{
		{"qs://test-bucket/obj", false, common.NewBucketDeleteTask, nil},
		{"qs://test-bucket/obj", true, common.NewRemoveBucketForceTask, nil},
		{"error", false, nil, removeBucketErr},
	}

	for _, v := range cases {
		pt := NewRemoveBucketTask(func(task *RemoveBucketTask) {
			task.SetForce(v.force)
			if v.expectErr != nil {
				task.TriggerFault(v.expectErr)
			}
		})

		assert.Equal(t,
			fmt.Sprintf("%v", v.expectedTodoFunc),
			fmt.Sprintf("%v", pt.Todo.NextTODO()))

		if v.expectErr != nil {
			assert.Equal(t, true, pt.ValidateFault())
			assert.Equal(t, true, errors.Is(pt.GetFault(), v.expectErr))
		} else {
			assert.Equal(t, false, pt.ValidateFault())
		}
	}
}

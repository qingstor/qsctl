package task

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yunify/qsctl/v2/pkg/types"
	"github.com/yunify/qsctl/v2/task/common"
)

func TestNewPresignTask(t *testing.T) {
	cases := []struct {
		input            string
		expectedTodoFunc types.TaskFunc
	}{
		{"qs://test-bucket/yyyyy", common.NewObjectPresignTask},
	}

	for _, v := range cases {
		pt := NewPresignTask(func(task *PresignTask) {})

		assert.Equal(t,
			fmt.Sprintf("%v", v.expectedTodoFunc),
			fmt.Sprintf("%v", pt.Todo.NextTODO()))
	}
}

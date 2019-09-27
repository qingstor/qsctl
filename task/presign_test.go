package task

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yunify/qsctl/v2/pkg/types"
)

func TestNewPresignTask(t *testing.T) {
	cases := []struct {
		input            string
		expectedTodoFunc types.TodoFunc
	}{
		{"qs://test-bucket/yyyyy", NewObjectPresignTask},
	}

	for _, v := range cases {
		pt := NewPresignTask(func(task *PresignTask) {
		})

		assert.Equal(t,
			fmt.Sprintf("%v", v.expectedTodoFunc),
			fmt.Sprintf("%v", pt.Todo.NextTODO()))
	}
}

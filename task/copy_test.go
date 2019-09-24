package task

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yunify/qsctl/v2/task/types"
	"github.com/yunify/qsctl/v2/task/utils"
)

func TestNewCopyTask(t *testing.T) {
	name, _, _ := utils.GenerateTestFile()
	defer os.Remove(name)

	cases := []struct {
		input1           string
		input2           string
		expectedTodoFunc types.TodoFunc
	}{
		{name, "qs://test-bucket/yyyyy", NewCopyFileTask},
		{"-", "qs://test-bucket/yyyyy", NewCopyStreamTask},
	}

	for _, v := range cases {
		pt := NewCopyTask(func(task *CopyTask) {
			err := utils.ParseInput(task, v.input1, v.input2)
			if err != nil {
				t.Fatal(err)
			}
		})

		assert.Equal(t,
			fmt.Sprintf("%v", v.expectedTodoFunc),
			fmt.Sprintf("%v", pt.Todo.NextTODO()))
	}

}

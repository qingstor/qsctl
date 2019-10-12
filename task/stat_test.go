package task

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yunify/qsctl/v2/pkg/types"
	"github.com/yunify/qsctl/v2/task/common"
	"github.com/yunify/qsctl/v2/utils"
)

func TestNewStatTask(t *testing.T) {
	cases := []struct {
		input            string
		expectedTodoFunc types.TodoFunc
		expectErr        error
	}{
		{"qs://test-bucket/obj", common.NewObjectStatTask, nil},
		// this test case is for PseudoDir, which will return error in the near future
		// {"qs://test-bucket/obj/", NewObjectStatTask, nil},
	}

	for _, v := range cases {
		pt := NewStatTask(func(task *StatTask) {
			_, _, _, err := utils.ParseKey(v.input)
			if err != nil {
				t.Fatal(err)
			}
		})

		assert.Equal(t,
			fmt.Sprintf("%v", v.expectedTodoFunc),
			fmt.Sprintf("%v", pt.Todo.NextTODO()))
	}
}

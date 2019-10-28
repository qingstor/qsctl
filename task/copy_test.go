package task

import (
	"fmt"
	"testing"

	typ "github.com/Xuanwo/storage/types"
	"github.com/stretchr/testify/assert"
	"github.com/yunify/qsctl/v2/pkg/types"
)

func TestNewCopyTask(t *testing.T) {
	cases := []struct {
		inputType        typ.ObjectType
		expectedTodoFunc types.TaskFunc
	}{
		{typ.ObjectTypeFile, NewCopyFileTask},
		{typ.ObjectTypeStream, NewCopyStreamTask},
	}

	for _, v := range cases {
		pt := NewCopyTask(func(task *CopyTask) {
			task.SetSourceType(v.inputType)
		})

		assert.Equal(t,
			fmt.Sprintf("%v", v.expectedTodoFunc),
			fmt.Sprintf("%v", pt.Todo.NextTODO()))
	}

}

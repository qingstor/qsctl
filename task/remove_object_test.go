package task

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yunify/qsctl/v2/pkg/types"
	"github.com/yunify/qsctl/v2/task/common"
)

func TestNewRemoveObjectTask(t *testing.T) {
	removeObjectErr := errors.New("remove-object-err")
	tests := []struct {
		name string
		fn   func(*RemoveObjectTask)
		want types.TodoFunc
		err  error
	}{
		{
			name: "obj",
			fn:   func(task *RemoveObjectTask) { task.SetRecursive(false) },
			want: common.NewObjectDeleteTask,
			err:  nil,
		},
		{
			name: "dir",
			fn:   func(task *RemoveObjectTask) { task.SetRecursive(true) },
			want: common.NewRemoveDirTask,
			err:  nil,
		},
		{
			name: "err",
			fn: func(task *RemoveObjectTask) {
				task.SetRecursive(true)
				task.SetFault(removeObjectErr)
			},
			want: nil,
			err:  removeObjectErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewRemoveObjectTask(tt.fn)
			assert.Equal(t,
				fmt.Sprintf("%v", tt.want),
				fmt.Sprintf("%v", got.NextTODO()))

			if tt.err != nil {
				assert.Equal(t, true, errors.Is(got.GetFault(), tt.err))
			} else {
				assert.Equal(t, false, got.ValidateFault())
			}
		})
	}
}

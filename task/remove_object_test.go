package task

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/yunify/qsctl/v2/pkg/types"
	"github.com/yunify/qsctl/v2/task/common"
)

func TestNewRemoveObjectTask(t *testing.T) {
	removeObjErr := errors.New("remove-obj-error")

	tests := []struct {
		name string
		fn   func(*RemoveObjectTask)
		want types.TodoFunc
		err  error
	}{
		{
			name: "remove obj",
			fn: func(task *RemoveObjectTask) {
				task.SetRecursive(false)
			},
			want: common.NewObjectDeleteTask,
			err:  nil,
		},
		{
			name: "remove dir",
			fn: func(task *RemoveObjectTask) {
				task.SetRecursive(true)
				task.SetKey(uuid.New().String())
			},
			want: common.NewRemoveDirTask,
			err:  nil,
		},
		{
			name: "got fault",
			fn: func(task *RemoveObjectTask) {
				task.TriggerFault(removeObjErr)
			},
			want: nil,
			err:  removeObjErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewRemoveObjectTask(tt.fn)
			if tt.err != nil {
				assert.Equal(t, true, got.ValidateFault())
				assert.Equal(t, true, errors.Is(got.GetFault(), tt.err))
			}
			assert.Equal(t,
				fmt.Sprintf("%v", tt.want),
				fmt.Sprintf("%v", got.NextTODO()))
		})
	}
}

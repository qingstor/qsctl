package task

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yunify/qsctl/v2/pkg/types"
	"github.com/yunify/qsctl/v2/task/common"
)

func TestNewRemoveObjectTask(t *testing.T) {
	type args struct {
		fn func(*RemoveObjectTask)
	}
	tests := []struct {
		name string
		args args
		want types.TodoFunc
	}{
		{name: "next", args: args{func(task *RemoveObjectTask) { task.SetRecursive(false) }}, want: common.NewObjectDeleteTask},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewRemoveObjectTask(tt.args.fn)
			assert.Equal(t,
				fmt.Sprintf("%v", tt.want),
				fmt.Sprintf("%v", got.NextTODO()))
		})
	}
}

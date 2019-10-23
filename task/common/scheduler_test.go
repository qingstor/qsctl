package common

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/yunify/qsctl/v2/pkg/types"
)

func TestDoneSchedulerTask_Run(t *testing.T) {
	id, doneErr := uuid.New().String(), errors.New("done-scheduler-error")
	cases := []struct {
		name string
		err  error
	}{
		{
			name: "ok",
			err:  nil,
		},
		{
			name: "error",
			err:  doneErr,
		},
	}

	for _, tt := range cases {
		x := &mockDoneSchedulerTask{}
		sche := types.NewMockScheduler(nil)
		x.SetScheduler(sche)
		x.SetID(id)
		if tt.err != nil {
			x.SetFault(tt.err)
		}

		sche.New(nil)

		task := NewDoneSchedulerTask(x)
		task.Run()

		if tt.err != nil {
			assert.Error(t, x.GetFault(), tt.name)
			assert.Equal(t, true, errors.Is(x.GetFault(), tt.err), tt.name)
		} else {
			assert.Equal(t, false, x.ValidateFault())
		}
	}

}

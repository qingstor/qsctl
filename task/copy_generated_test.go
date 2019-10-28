// Code generated by go generate; DO NOT EDIT.
package task

import (
	"errors"
	"testing"

	"github.com/Xuanwo/navvy"
	"github.com/stretchr/testify/assert"

	"github.com/yunify/qsctl/v2/pkg/types"
)

var _ navvy.Pool
var _ types.Pool

func TestNewCopyTask(t *testing.T) {
	m := &mockCopyTask{}
	task := NewCopyTask(m)
	assert.NotNil(t, task)
}

func TestCopyTask_GeneratedRun(t *testing.T) {
	cases := []struct {
		name     string
		hasFault bool
		hasCall  bool
		gotCall  bool
	}{
		{
			"has fault",
			true,
			false,
			false,
		},
		{
			"no fault",
			false,
			true,
			false,
		},
	}

	for _, v := range cases {
		t.Run(v.name, func(t *testing.T) {
			pool := navvy.NewPool(10)
			task := &CopyTask{}
			task.SetPool(pool)

			err := errors.New("test error")
			if v.hasFault {
				task.SetFault(err)
			}
			task.GetScheduler.Sync(func(todoist types.TaskFunc) navvy.Task {
				x := utils.NewCallbackTask(func() {
					v.gotCall = true
				})
				return x
			}, task)

			task.Run()
			pool.Wait()

			assert.Equal(t, v.hasCall, v.gotCall)
		})
	}
}

func TestCopyTask_TriggerFault(t *testing.T) {
	m := &mockCopyTask{}
	task := &CopyTask{m}
	err := errors.New("test error")
	task.TriggerFault(err)
	assert.True(t, task.copyTaskRequirement.ValidateFault())
}

func TestMockCopyTask_Run(t *testing.T) {
	task := &mockCopyTask{}
	assert.Panics(t, func() {
		task.Run()
	})
}
func TestCopyTask_Wait(t *testing.T) {
	pool := navvy.NewPool(10)
	task := &CopyTask{}
	{
		assert.Panics(t, func() {
			task.Wait()
		})
	}
	{
		task.SetPool(pool)
		assert.NotPanics(t, func() {
			task.Wait()
		})
	}
}

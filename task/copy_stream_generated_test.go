// Code generated by go generate; DO NOT EDIT.
package task

import (
	"errors"
	"testing"

	"github.com/Xuanwo/navvy"
	"github.com/stretchr/testify/assert"

	"github.com/yunify/qsctl/v2/pkg/types"
	"github.com/yunify/qsctl/v2/utils"
)

var _ navvy.Pool
var _ types.Pool
var _ = utils.SubmitNextTask

func TestCopyPartialStreamTask_GeneratedRun(t *testing.T) {
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
			m := &mockCopyPartialStreamTask{}
			m.SetPool(pool)
			task := &CopyPartialStreamTask{copyPartialStreamTaskRequirement: m}

			err := errors.New("test error")
			if v.hasFault {
				task.SetFault(err)
			}
			task.AddTODOs(func(todoist types.Todoist) navvy.Task {
				x := utils.NewCallbackTask(func() {
					v.gotCall = true
				})
				return x
			})

			task.Run()
			pool.Wait()

			assert.Equal(t, v.hasCall, v.gotCall)
		})
	}
}

func TestCopyPartialStreamTask_TriggerFault(t *testing.T) {
	err := errors.New("trigger fault")
	x := &CopyPartialStreamTask{}
	x.TriggerFault(err)

	assert.Equal(t, true, x.ValidateFault())
	assert.Equal(t, true, errors.Is(x.GetFault(), err))
}

func TestCopyStreamTask_GeneratedRun(t *testing.T) {
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
			m := &mockCopyStreamTask{}
			m.SetPool(pool)
			task := &CopyStreamTask{copyStreamTaskRequirement: m}

			err := errors.New("test error")
			if v.hasFault {
				task.SetFault(err)
			}
			task.AddTODOs(func(todoist types.Todoist) navvy.Task {
				x := utils.NewCallbackTask(func() {
					v.gotCall = true
				})
				return x
			})

			task.Run()
			pool.Wait()

			assert.Equal(t, v.hasCall, v.gotCall)
		})
	}
}

func TestCopyStreamTask_TriggerFault(t *testing.T) {
	err := errors.New("trigger fault")
	x := &CopyStreamTask{}
	x.TriggerFault(err)

	assert.Equal(t, true, x.ValidateFault())
	assert.Equal(t, true, errors.Is(x.GetFault(), err))
}

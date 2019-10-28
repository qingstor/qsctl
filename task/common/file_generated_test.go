// Code generated by go generate; DO NOT EDIT.
package common

import (
	"errors"
	"testing"

	"github.com/Xuanwo/navvy"
	"github.com/stretchr/testify/assert"

	"github.com/yunify/qsctl/v2/pkg/types"
)

var _ navvy.Pool
var _ types.Pool

func TestNewFileUploadTask(t *testing.T) {
	m := &mockFileUploadTask{}
	task := NewFileUploadTask(m)
	assert.NotNil(t, task)
}

func TestFileUploadTask_GeneratedRun(t *testing.T) {
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

			m := &mockFileUploadTask{}
			m.SetPool(pool)
			task := &FileUploadTask{fileUploadTaskRequirement: m}

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

func TestFileUploadTask_TriggerFault(t *testing.T) {
	m := &mockFileUploadTask{}
	task := &FileUploadTask{m}
	err := errors.New("test error")
	task.TriggerFault(err)
	assert.True(t, task.fileUploadTaskRequirement.ValidateFault())
}

func TestMockFileUploadTask_Run(t *testing.T) {
	task := &mockFileUploadTask{}
	assert.Panics(t, func() {
		task.Run()
	})
}

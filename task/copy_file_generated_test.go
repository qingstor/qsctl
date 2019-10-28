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

func TestNewCopyFileTask(t *testing.T) {
	m := &mockCopyFileTask{}
	task := NewCopyFileTask(m)
	assert.NotNil(t, task)
}

func TestCopyFileTask_GeneratedRun(t *testing.T) {
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
			m := &mockCopyFileTask{}
			m.SetPool(pool)
			task := &CopyFileTask{copyFileTaskRequirement: m}

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

func TestCopyFileTask_TriggerFault(t *testing.T) {
	m := &mockCopyFileTask{}
	task := &CopyFileTask{m}
	err := errors.New("test error")
	task.TriggerFault(err)
	assert.True(t, task.copyFileTaskRequirement.ValidateFault())
}

func TestMockCopyFileTask_Run(t *testing.T) {
	task := &mockCopyFileTask{}
	assert.Panics(t, func() {
		task.Run()
	})
}

func TestNewCopyLargeFileTask(t *testing.T) {
	m := &mockCopyLargeFileTask{}
	task := NewCopyLargeFileTask(m)
	assert.NotNil(t, task)
}

func TestCopyLargeFileTask_GeneratedRun(t *testing.T) {
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
			m := &mockCopyLargeFileTask{}
			m.SetPool(pool)
			task := &CopyLargeFileTask{copyLargeFileTaskRequirement: m}

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

func TestCopyLargeFileTask_TriggerFault(t *testing.T) {
	m := &mockCopyLargeFileTask{}
	task := &CopyLargeFileTask{m}
	err := errors.New("test error")
	task.TriggerFault(err)
	assert.True(t, task.copyLargeFileTaskRequirement.ValidateFault())
}

func TestMockCopyLargeFileTask_Run(t *testing.T) {
	task := &mockCopyLargeFileTask{}
	assert.Panics(t, func() {
		task.Run()
	})
}

func TestNewCopyPartialFileTask(t *testing.T) {
	m := &mockCopyPartialFileTask{}
	task := NewCopyPartialFileTask(m)
	assert.NotNil(t, task)
}

func TestCopyPartialFileTask_GeneratedRun(t *testing.T) {
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
			m := &mockCopyPartialFileTask{}
			m.SetPool(pool)
			task := &CopyPartialFileTask{copyPartialFileTaskRequirement: m}

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

func TestCopyPartialFileTask_TriggerFault(t *testing.T) {
	m := &mockCopyPartialFileTask{}
	task := &CopyPartialFileTask{m}
	err := errors.New("test error")
	task.TriggerFault(err)
	assert.True(t, task.copyPartialFileTaskRequirement.ValidateFault())
}

func TestMockCopyPartialFileTask_Run(t *testing.T) {
	task := &mockCopyPartialFileTask{}
	assert.Panics(t, func() {
		task.Run()
	})
}

func TestNewCopySmallFileTask(t *testing.T) {
	m := &mockCopySmallFileTask{}
	task := NewCopySmallFileTask(m)
	assert.NotNil(t, task)
}

func TestCopySmallFileTask_GeneratedRun(t *testing.T) {
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
			m := &mockCopySmallFileTask{}
			m.SetPool(pool)
			task := &CopySmallFileTask{copySmallFileTaskRequirement: m}

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

func TestCopySmallFileTask_TriggerFault(t *testing.T) {
	m := &mockCopySmallFileTask{}
	task := &CopySmallFileTask{m}
	err := errors.New("test error")
	task.TriggerFault(err)
	assert.True(t, task.copySmallFileTaskRequirement.ValidateFault())
}

func TestMockCopySmallFileTask_Run(t *testing.T) {
	task := &mockCopySmallFileTask{}
	assert.Panics(t, func() {
		task.Run()
	})
}

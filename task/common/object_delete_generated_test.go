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

func TestNewObjectDeleteTask(t *testing.T) {
	m := &mockObjectDeleteTask{}
	task := NewObjectDeleteTask(m)
	assert.NotNil(t, task)
}

func TestObjectDeleteTask_GeneratedRun(t *testing.T) {
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
			task := &ObjectDeleteTask{}
			task.SetPool(pool)

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

func TestObjectDeleteTask_TriggerFault(t *testing.T) {
	m := &mockObjectDeleteTask{}
	task := &ObjectDeleteTask{m}
	err := errors.New("test error")
	task.TriggerFault(err)
	assert.True(t, task.objectDeleteTaskRequirement.ValidateFault())
}

func TestMockObjectDeleteTask_Run(t *testing.T) {
	task := &mockObjectDeleteTask{}
	assert.Panics(t, func() {
		task.Run()
	})
}
func TestObjectDeleteTask_Wait(t *testing.T) {
	pool := navvy.NewPool(10)
	task := &ObjectDeleteTask{}
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

func TestNewObjectDeleteIterateTask(t *testing.T) {
	m := &mockObjectDeleteIterateTask{}
	task := NewObjectDeleteIterateTask(m)
	assert.NotNil(t, task)
}

func TestObjectDeleteIterateTask_GeneratedRun(t *testing.T) {
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
			task := &ObjectDeleteIterateTask{}
			task.SetPool(pool)

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

func TestObjectDeleteIterateTask_TriggerFault(t *testing.T) {
	m := &mockObjectDeleteIterateTask{}
	task := &ObjectDeleteIterateTask{m}
	err := errors.New("test error")
	task.TriggerFault(err)
	assert.True(t, task.objectDeleteIterateTaskRequirement.ValidateFault())
}

func TestMockObjectDeleteIterateTask_Run(t *testing.T) {
	task := &mockObjectDeleteIterateTask{}
	assert.Panics(t, func() {
		task.Run()
	})
}
func TestObjectDeleteIterateTask_Wait(t *testing.T) {
	pool := navvy.NewPool(10)
	task := &ObjectDeleteIterateTask{}
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

func TestNewObjectDeleteScheduledTask(t *testing.T) {
	m := &mockObjectDeleteScheduledTask{}
	task := NewObjectDeleteScheduledTask(m)
	assert.NotNil(t, task)
}

func TestObjectDeleteScheduledTask_GeneratedRun(t *testing.T) {
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
			m := &mockObjectDeleteScheduledTask{}
			m.SetPool(pool)
			task := &ObjectDeleteScheduledTask{objectDeleteScheduledTaskRequirement: m}

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

func TestObjectDeleteScheduledTask_TriggerFault(t *testing.T) {
	m := &mockObjectDeleteScheduledTask{}
	task := &ObjectDeleteScheduledTask{m}
	err := errors.New("test error")
	task.TriggerFault(err)
	assert.True(t, task.objectDeleteScheduledTaskRequirement.ValidateFault())
}

func TestMockObjectDeleteScheduledTask_Run(t *testing.T) {
	task := &mockObjectDeleteScheduledTask{}
	assert.Panics(t, func() {
		task.Run()
	})
}

func TestNewRemoveDirTask(t *testing.T) {
	m := &mockRemoveDirTask{}
	task := NewRemoveDirTask(m)
	assert.NotNil(t, task)
}

func TestRemoveDirTask_GeneratedRun(t *testing.T) {
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
			m := &mockRemoveDirTask{}
			m.SetPool(pool)
			task := &RemoveDirTask{removeDirTaskRequirement: m}

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

func TestRemoveDirTask_TriggerFault(t *testing.T) {
	m := &mockRemoveDirTask{}
	task := &RemoveDirTask{m}
	err := errors.New("test error")
	task.TriggerFault(err)
	assert.True(t, task.removeDirTaskRequirement.ValidateFault())
}

func TestMockRemoveDirTask_Run(t *testing.T) {
	task := &mockRemoveDirTask{}
	assert.Panics(t, func() {
		task.Run()
	})
}

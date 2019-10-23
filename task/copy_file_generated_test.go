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

func TestCopyFileTask_TriggerFault(t *testing.T) {
	err := errors.New("trigger fault")
	x := &CopyFileTask{}
	x.TriggerFault(err)

	assert.Equal(t, true, x.ValidateFault())
	assert.Equal(t, true, errors.Is(x.GetFault(), err))
}

func TestMockCopyFileTask_Run(t *testing.T) {
	task := &mockCopyFileTask{}
	assert.Panics(t, func() {
		task.Run()
	})
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

func TestCopyLargeFileTask_TriggerFault(t *testing.T) {
	err := errors.New("trigger fault")
	x := &CopyLargeFileTask{}
	x.TriggerFault(err)

	assert.Equal(t, true, x.ValidateFault())
	assert.Equal(t, true, errors.Is(x.GetFault(), err))
}

func TestMockCopyLargeFileTask_Run(t *testing.T) {
	task := &mockCopyLargeFileTask{}
	assert.Panics(t, func() {
		task.Run()
	})
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

func TestCopyPartialFileTask_TriggerFault(t *testing.T) {
	err := errors.New("trigger fault")
	x := &CopyPartialFileTask{}
	x.TriggerFault(err)

	assert.Equal(t, true, x.ValidateFault())
	assert.Equal(t, true, errors.Is(x.GetFault(), err))
}

func TestMockCopyPartialFileTask_Run(t *testing.T) {
	task := &mockCopyPartialFileTask{}
	assert.Panics(t, func() {
		task.Run()
	})
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

func TestCopySmallFileTask_TriggerFault(t *testing.T) {
	err := errors.New("trigger fault")
	x := &CopySmallFileTask{}
	x.TriggerFault(err)

	assert.Equal(t, true, x.ValidateFault())
	assert.Equal(t, true, errors.Is(x.GetFault(), err))
}

func TestMockCopySmallFileTask_Run(t *testing.T) {
	task := &mockCopySmallFileTask{}
	assert.Panics(t, func() {
		task.Run()
	})
}

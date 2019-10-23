// Code generated by go generate; DO NOT EDIT.
package common

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

func TestNewFileMD5SumTask(t *testing.T) {
	m := &mockFileMD5SumTask{}
	task := NewFileMD5SumTask(m)
	assert.NotNil(t, task)
}

func TestFileMD5SumTask_TriggerFault(t *testing.T) {
	m := &mockFileMD5SumTask{}
	task := &FileMD5SumTask{m}
	err := errors.New("test error")
	task.TriggerFault(err)
	assert.True(t, task.fileMD5SumTaskRequirement.ValidateFault())
}

func TestMockFileMD5SumTask_Run(t *testing.T) {
	task := &mockFileMD5SumTask{}
	assert.Panics(t, func() {
		task.Run()
	})
}

func TestNewStreamMD5SumTask(t *testing.T) {
	m := &mockStreamMD5SumTask{}
	task := NewStreamMD5SumTask(m)
	assert.NotNil(t, task)
}

func TestStreamMD5SumTask_TriggerFault(t *testing.T) {
	m := &mockStreamMD5SumTask{}
	task := &StreamMD5SumTask{m}
	err := errors.New("test error")
	task.TriggerFault(err)
	assert.True(t, task.streamMD5SumTaskRequirement.ValidateFault())
}

func TestMockStreamMD5SumTask_Run(t *testing.T) {
	task := &mockStreamMD5SumTask{}
	assert.Panics(t, func() {
		task.Run()
	})
}

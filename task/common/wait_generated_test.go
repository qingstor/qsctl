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

func TestNewWaitTask(t *testing.T) {
	m := &mockWaitTask{}
	task := NewWaitTask(m)
	assert.NotNil(t, task)
}

func TestWaitTask_TriggerFault(t *testing.T) {
	m := &mockWaitTask{}
	task := &WaitTask{m}
	err := errors.New("test error")
	task.TriggerFault(err)
	assert.True(t, task.waitTaskRequirement.ValidateFault())
}
